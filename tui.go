package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	tui "github.com/marcusolsson/tui-go"
	"github.com/sirupsen/logrus"
)

type post struct {
	ID       int
	username string
	message  string
	time     string
}

// ByID implements sort.Interface for []post based on
// the ID field.
type ByID []post

func (id ByID) Len() int           { return len(id) }
func (id ByID) Swap(i, j int)      { id[i], id[j] = id[j], id[i] }
func (id ByID) Less(i, j int) bool { return id[i].ID < id[j].ID }

func (p *post) messageToPost(message *discordgo.Message) error {
	p.username = message.Author.String()
	p.message = message.ContentWithMentionsReplaced()
	parsedTime, err := message.Timestamp.Parse()
	if err != nil {
		return err
	}
	p.time = parsedTime.Format("15:04")
	return nil
}

// Channel Channel being displayed by TUI
var (
	channel       *discordgo.Channel
	posts         []post
	tuiDS         *discordgo.Session
	history       = tui.NewVBox()
	historyScroll = tui.NewScrollArea(history)
	historyBox    = tui.NewVBox(historyScroll)
	input         = tui.NewEntry()
	inputBox      = tui.NewHBox(input)
	chat          = tui.NewVBox(historyBox, inputBox)
	root          = tui.NewHBox(chat)
	ui            tui.UI
	postCount     = 0
	lastAck       = new(discordgo.Ack)
)

const (
	tuiHelp = `TUI for sending and receiving Discord DMs`
	// messagesFetchLimit is the limit that discordgo returns of old messages. Max is 100
	messagesFetchLimit = 100
)

func init() {
	localUI, err := tui.New(root)
	if err != nil {
		panic(err, tuiDS)
	}
	ui = localUI
}

func (cmd *tuiCommand) Name() string      { return "tui" }
func (cmd *tuiCommand) Args() string      { return "[OPTIONS]" }
func (cmd *tuiCommand) ShortHelp() string { return tuiHelp }
func (cmd *tuiCommand) LongHelp() string  { return tuiHelp }
func (cmd *tuiCommand) Hidden() bool      { return false }

func (cmd *tuiCommand) Register(fs *flag.FlagSet) {
	fs.IntVar(&cmd.cSel, "c", -1, "specify the channel to start the tui in")
}

type tuiCommand struct {
	cSel int
}

func (cmd *tuiCommand) Run(ctx context.Context, args []string) error {
	ds := getContextValue(ctx, discordSessionKey)
	StartTUI(cmd.cSel, ds)
	return nil
}

// Error Handler
func panic(err error, ds *discordgo.Session) {
	fmt.Printf("%v\n", err)
	tuiDS.Close()
	os.Exit(1)
}

// StartTUI Start the TUI display
func StartTUI(cSel int, ds *discordgo.Session) {
	// sidebar := tui.NewVBox(
	// 	tui.NewLabel("CHANNELS"),
	// 	tui.NewLabel("general"),
	// 	tui.NewLabel("random"),
	// 	tui.NewLabel(""),
	// 	tui.NewLabel("DIRECT MESSAGES"),
	// 	tui.NewLabel("slackbot"),
	// 	tui.NewSpacer(),
	// )
	// sidebar.SetBorder(true)

	// Get a list of DM channels
	channels, err := tuiDS.UserChannels()
	logrus.Debugf("Retrieved Channels\n%v\n", spew.Sdump(channels))
	if err != nil {
		panic(err, tuiDS)
	}

	channelSel := cSel
	if channelSel <= -1 {
		// Display available channels
		fmt.Println("Available Private Channels:")
		for i, chann := range channels {
			logrus.Debugf("Channel %d\n%v\n", i, spew.Sdump(chann))

			// Switch of supported channel types
			switch chanType := chann.Type; chanType {
			case 1: // Direct Messages
				var flatRecipients string
				recipients := chann.Recipients
				logrus.Debugf("Channel %d recipients\n%v\n", i, spew.Sdump(recipients))
				for _, recipient := range recipients {
					flatRecipients = fmt.Sprintf("%s %s", flatRecipients, recipient.Username)
				}
				fmt.Printf("\t%d) DM to %s\n", i, strings.TrimSpace(flatRecipients))
			default:
				fmt.Println("No available channels")
				os.Exit(0)
			}
		}

		// Get channel selection
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Select a channel to switch to: ")
		channelSelS, err := reader.ReadString('\n')
		if err != nil {
			panic(err, tuiDS)
		}
		channelSel, err = strconv.Atoi(strings.TrimSpace(channelSelS))
		if err != nil {
			panic(err, tuiDS)
		}
	}

	// Set Channel to selected channel
	if (len(channels) - 1) < channelSel {
		err := fmt.Errorf("Channel selection outside of available selection range\n%d < %d", (len(channels) - 1), channelSel)
		panic(err, tuiDS)
	}
	channel = channels[channelSel]
	logrus.Debugf("Channel selected %d\n%v\n", channelSel, spew.Sdump(channel))

	// Get Channel messages
	messages, err := tuiDS.ChannelMessages(channel.ID, messagesFetchLimit, "", "", "")
	if err != nil {
		panic(err, tuiDS)
	}
	logrus.Debugf("Channel messages\n%v\n", spew.Sdump(messages))
	// spew.Dump(messages[0])

	// Convert messages to posts
	posts, err := convertToTUIPosts(messages)
	if err != nil {
		panic(err, tuiDS)
	}

	// // sort posts by id
	// sort.Sort(ByID(posts))
	// logrus.Debugf("# of Channel messages\n%d\n", len(posts))

	// history defined above as the location to display messages
	// history = tui.NewVBox()

	// Add posts to history box
	addPostsToDisplay(posts)

	// historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	// historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	// input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	// inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	// chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		_, err := tuiDS.ChannelMessageSend(channel.ID, e.Text())
		if err != nil {
			panic(err, tuiDS)
		}
		// p, err := convertToTUIPost(message)
		// if err != nil {
		// 	panic(err, tuiDS)
		// }
		// addPostToDisplay(p)
		input.SetText("")
	})

	// root := tui.NewHBox(chat)

	// ui, err := tui.New(root)
	// if err != nil {
	// 	panic(err, tuiDS)
	// }

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	// Register the messageCreate func as a callback for MessageCreate events.
	tuiDS.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = tuiDS.Open()
	if err != nil {
		err = fmt.Errorf("error opening connection, %v", err)
		panic(err, tuiDS)
	}

	if err := ui.Run(); err != nil {
		panic(err, tuiDS)
	}
	// Cleanly close down the Discord session.
	tuiDS.Close()
}

// Add posts to the history box (defined above)
func addPostsToDisplay(ps []post) {
	// sort posts by id
	sort.Sort(sort.Reverse(ByID(ps)))
	logrus.Debugf("# of Channel messages\n%d\n", len(ps))
	for _, p := range ps {
		p.message = strconv.QuoteToASCII(p.message)
		history.Append(tui.NewHBox(
			tui.NewLabel(p.time),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", p.username))),
			tui.NewLabel(p.message),
			tui.NewSpacer(),
		))
	}
}

func addPostToDisplay(p post) {
	p.message = strconv.QuoteToASCII(p.message)
	ui.Update(func() {
		history.Append(tui.NewHBox(
			tui.NewLabel(p.time),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", p.username))),
			tui.NewLabel(p.message),
			tui.NewSpacer(),
		))
	})
}

// Convert discordgo messages to posts
func convertToTUIPosts(messages []*discordgo.Message) ([]post, error) {
	var ps []post
	for _, message := range messages {
		p, err := convertToTUIPost(message)
		if err != nil {
			return ps, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

// Convert discordgo messages to posts
func convertToTUIPost(message *discordgo.Message) (post, error) {
	p := post{ID: postCount}
	if err := p.messageToPost(message); err != nil {
		return p, err
	}
	postCount++
	return p, nil
}

// New message event handler
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message.ChannelID == channel.ID {
		newAck, err := s.ChannelMessageAck(m.Message.ChannelID, m.Message.ID, lastAck.Token)
		if err != nil {
			panic(err, tuiDS)
		}
		lastAck = newAck
		p, err := convertToTUIPost(m.Message)
		if err != nil {
			panic(err, tuiDS)
		}
		addPostToDisplay(p)
	}
}
