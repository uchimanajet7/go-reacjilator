package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

// SlackClient is
type SlackClient struct {
	client            *slack.Client
	verificationToken string
	channelID         string
}

type slackMsg struct {
	text       string
	ts         string
	channel    string
	reaction   string
	translated string
	source     string
	target     string
}

// automatically generated using the following
// https://mholt.github.io/json-to-go/
type slackEvent struct {
	// only url_verification event
	// https://api.slack.com/events/url_verification
	Challenge string `json:"challenge"`

	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	APIAppID string `json:"api_app_id"`
	Event    struct {
		Type string `json:"type"`
		User string `json:"user"`
		Item struct {
			Type    string `json:"type"`
			Channel string `json:"channel"`
			Ts      string `json:"ts"`
		} `json:"item"`
		Reaction string `json:"reaction"`
		ItemUser string `json:"item_user"`
		EventTs  string `json:"event_ts"`
	} `json:"event"`
	Type        string   `json:"type"`
	EventID     string   `json:"event_id"`
	EventTime   int      `json:"event_time"`
	AuthedUsers []string `json:"authed_users"`
}

// https://api.slack.com/events/url_verification
/*
{
    "token": "Jhj5dZrVaK7ZwHHjRyZWjbDl",
    "challenge": "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P",
    "type": "url_verification"
}
*/

// https://api.slack.com/events-api#receiving_events
// https://api.slack.com/events/reaction_added
/*
{
    "token": "Jhj5dZrVaK7ZwHHjRyZWjbDl",
    "team_id": "T0123086H",
    "api_app_id": "A1231KF9R",
    "event": {
        "type": "reaction_added",
        "user": "U0G9QF9C6",
        "item": {
            "type": "message",
            "channel": "C0G9QF9GZ",
            "ts": "1518504389.000166"
        },
        "reaction": "eyes",
        "item_user": "U0G9QF9C6",
        "event_ts": "1518507482.000119"
    },
    "type": "event_callback",
    "event_id": "Ev97FS5N0Y",
    "event_time": 1518507482,
    "authed_users": [
        "U0G9QF9C6"
    ]
}
*/

// flag emoji to language
var flagMap = map[string]string{
	"flag-ac": "English",
	"flag-ad": "Catalan",
	"flag-ae": "Arabic",
	"flag-af": "Pashto",
	"flag-ag": "English",
	"flag-ai": "English",
	"flag-al": "Albanian",
	"flag-am": "Armenian",
	"flag-ao": "Portuguese",
	"flag-ar": "Spanish",
	"flag-as": "English",
	"flag-at": "German",
	"flag-au": "English",
	"flag-aw": "Dutch",
	"flag-ax": "Swedish",
	"flag-az": "Spanish",
	"flag-ba": "Bosnian",
	"flag-bb": "English",
	"flag-bd": "Bengali",
	"flag-be": "Dutch",
	"flag-bf": "French",
	"flag-bg": "Bulgarian",
	"flag-bh": "Arabic",
	"flag-bi": "French",
	"flag-bj": "French",
	"flag-bl": "French",
	"flag-bn": "English",
	"flag-bm": "Malay",
	"flag-bo": "Spanish",
	"flag-bq": "Dutch",
	"flag-br": "Portuguese",
	"flag-bs": "English",
	"flag-bt": "Dzongkha",
	"flag-bv": "Norwegian",
	"flag-bw": "English",
	"flag-by": "Belarusian",
	"flag-bz": "English",
	"flag-ca": "English",
	"flag-cc": "Malay",
	"flag-cd": "French",
	"flag-cf": "French",
	"flag-cg": "French",
	"flag-ch": "German",
	"flag-ci": "French",
	"flag-ck": "English",
	"flag-cl": "Spanish",
	"flag-cm": "French",
	"flag-cn": "Chinese Simplified",
	"flag-co": "Spanish",
	"flag-cp": "French",
	"flag-cr": "Spanish",
	"flag-cu": "Spanish",
	"flag-cv": "Portuguese",
	"flag-cw": "Dutch",
	"flag-cx": "English",
	"flag-cy": "Greek",
	"flag-cz": "Czech",
	"flag-de": "German",
	"flag-dg": "English",
	"flag-dj": "French",
	"flag-dk": "Danish",
	"flag-dm": "English",
	"flag-do": "Spanish",
	"flag-dz": "Arabic",
	"flag-ea": "Spanish",
	"flag-ec": "Spanish",
	"flag-ee": "Estonian",
	"flag-eg": "Arabic",
	"flag-eh": "Arabic",
	"flag-er": "Arabic",
	"flag-es": "Spanish",
	"flag-et": "Oromo",
	"flag-fi": "Finnish",
	"flag-fj": "English",
	"flag-fk": "English",
	"flag-fm": "English",
	"flag-fr": "French",
	"flag-ga": "French",
	"flag-gb": "English",
	"flag-gd": "English",
	"flag-ge": "Georgian",
	"flag-gf": "French",
	"flag-gg": "English",
	"flag-gh": "English",
	"flag-gi": "English",
	"flag-gl": "Danish",
	"flag-gm": "English",
	"flag-gn": "French",
	"flag-gp": "French",
	"flag-gq": "Spanish",
	"flag-gr": "Greek",
	"flag-gs": "English",
	"flag-gt": "Spanish",
	"flag-gu": "English",
	"flag-gw": "Portuguese",
	"flag-gy": "English",
	"flag-hk": "Chinese Traditional",
	"flag-hn": "Spanish",
	"flag-hr": "Croatian",
	"flag-ht": "Haitian Creole",
	"flag-hu": "Hungarian",
	"flag-ic": "Spanish",
	"flag-id": "Indonesian",
	"flag-ie": "Irish",
	"flag-il": "Hebrew",
	"flag-im": "English",
	"flag-in": "Hindi",
	"flag-io": "English",
	"flag-iq": "Arabic",
	"flag-ir": "Persian",
	"flag-is": "Icelandic",
	"flag-it": "Italian",
	"flag-je": "English",
	"flag-jm": "English",
	"flag-jo": "Arabic",
	"flag-jp": "Japanese",
	"flag-ke": "English",
	"flag-kg": "Kyrgyz",
	"flag-kh": "Khmer",
	"flag-ki": "English",
	"flag-kn": "English",
	"flag-kp": "Korean",
	"flag-kr": "Korean",
	"flag-kw": "Arabic",
	"flag-ky": "English",
	"flag-kz": "Kazakh",
	"flag-la": "Lao",
	"flag-lb": "Arabic",
	"flag-lc": "English",
	"flag-li": "German",
	"flag-lk": "Sinhala",
	"flag-lr": "English",
	"flag-ls": "Sesotho",
	"flag-lt": "Lithuanian",
	"flag-lu": "Luxembourgish",
	"flag-lv": "Latvian",
	"flag-ly": "Arabic",
	"flag-ma": "Arabic",
	"flag-mc": "French",
	"flag-md": "Romanian",
	"flag-mg": "Malagasy",
	"flag-mh": "Marshallese",
	"flag-mk": "Macedonian",
	"flag-ml": "French",
	"flag-mm": "Burmese",
	"flag-mn": "Mongolian",
	"flag-mo": "Chinese Traditional",
	"flag-mp": "English",
	"flag-mq": "French",
	"flag-mr": "Arabic",
	"flag-ms": "English",
	"flag-mt": "Maltese",
	"flag-mu": "English",
	"flag-mv": "Dhivehi",
	"flag-mw": "English",
	"flag-mx": "Spanish",
	"flag-my": "Malay",
	"flag-mz": "Portuguese",
	"flag-na": "English",
	"flag-nc": "French",
	"flag-ne": "French",
	"flag-nf": "English",
	"flag-ng": "English",
	"flag-ni": "Spanish",
	"flag-nl": "Dutch",
	"flag-no": "Norwegian",
	"flag-np": "Nepali",
	"flag-nr": "Nauru",
	"flag-nu": "Niuean",
	"flag-nz": "English",
	"flag-om": "Arabic",
	"flag-pa": "Spanish",
	"flag-pe": "Spanish",
	"flag-pf": "French",
	"flag-pg": "English",
	"flag-ph": "Tagalog",
	"flag-pk": "Urdu",
	"flag-pl": "Polish",
	"flag-pm": "French",
	"flag-pn": "English",
	"flag-pr": "Spanish",
	"flag-ps": "Arabic",
	"flag-pt": "Portuguese",
	"flag-pw": "English",
	"flag-py": "Spanish",
	"flag-qa": "Arabic",
	"flag-re": "French",
	"flag-ro": "Romanian",
	"flag-rs": "Serbian",
	"flag-ru": "Russian",
	"flag-rw": "Kinyarwanda",
	"flag-sa": "Arabic",
	"flag-sb": "English",
	"flag-sc": "English",
	"flag-sd": "Arabic",
	"flag-se": "Swedish",
	"flag-sg": "English",
	"flag-sh": "English",
	"flag-si": "Slovenian",
	"flag-sj": "Norwegian",
	"flag-sk": "Slovak",
	"flag-sl": "English",
	"flag-sm": "Italian",
	"flag-sn": "French",
	"flag-so": "Somali",
	"flag-sr": "Dutch",
	"flag-ss": "English",
	"flag-st": "Portuguese",
	"flag-sv": "Spanish",
	"flag-sx": "Dutch",
	"flag-sw": "Arabic",
	"flag-sz": "Swati",
	"flag-ta": "English",
	"flag-tc": "English",
	"flag-td": "French",
	"flag-tf": "French",
	"flag-tg": "French",
	"flag-th": "Thai",
	"flag-tj": "Tajik",
	"flag-tk": "Tokelau",
	"flag-tl": "Tetum",
	"flag-tm": "Turkmen",
	"flag-tn": "Arabic",
	"flag-tr": "Turkish",
	"flag-tt": "English",
	"flag-tv": "Tuvalua",
	"flag-tw": "Chinese Traditional",
	"flag-tz": "Swahili",
	"flag-ua": "Ukrainian",
	"flag-ug": "English",
	"flag-um": "English",
	"flag-us": "English",
	"flag-uy": "Spanish",
	"flag-uz": "Uzbek",
	"flag-va": "Italian",
	"flag-vc": "English",
	"flag-ve": "Spanish",
	"flag-vg": "English",
	"flag-vi": "English",
	"flag-vn": "Vietnamese",
	"flag-vu": "English",
	"flag-wf": "French",
	"flag-ws": "Samoan",
	"flag-xk": "Albanian",
	"flag-ye": "Arabic",
	"flag-yt": "French",
	"flag-za": "Afrikaans",
	"flag-zm": "English",
	"flag-zw": "English",
	"flag-to": "",
	"flag-me": "",
	"flag-km": "",
	"flag-hm": "",
	"flag-mf": "Saint Martin",
	"flag-fo": "Faroe Islands",
	"flag-eu": "EU",
	"flag-aq": "Antarctica",
}

func (c *SlackClient) handleEvent(data string) (string, error) {
	var se slackEvent
	if err := json.Unmarshal([]byte(data), &se); err != nil {
		log.Println("[Error] JSON unmarshal error:", err)
		return "", err
	}

	// check verification token
	if se.Token != c.verificationToken {
		log.Println("[Error] slack verification token do not match error: ", se.Token)
		return "", errors.New("slack verification token do not match!!/n")
	}

	// url verification
	if se.Type == "url_verification" {
		log.Println("[Accepted] url_verification event")
		return fmt.Sprintf(`{"challenge": %s}`, se.Challenge), nil
	}

	if se.Event.Type != "reaction_added" {
		log.Println("[Rejected] slack event type do not 'reaction_added': ", se.Event.Type)
		return "", nil
	}

	// filter the channel?
	if c.channelID != "" {
		if c.channelID != se.Event.Item.Channel {
			log.Println("[Rejected] slack channel ID do not match: ", se.Event.Item.Channel)
			return "", nil
		}
	}

	// determine the language from the flag emoji
	reactionText := se.Event.Reaction
	if !strings.HasPrefix(reactionText, "flag-") {
		reactionText = "flag-" + reactionText
	}
	targetCode := GetLanguageCode(flagMap[reactionText])
	if targetCode == "" {
		log.Println("[Rejected] it does not correspond to that emoji reaction: ", se.Event.Reaction)
		return "", nil
	}

	// get slack message
	msg, err := c.getMessage(se.Event.Item.Channel, se.Event.Item.Ts)
	if err != nil {
		log.Println("[Error] failed to get slack messages: ", err)
		return "", errors.New("failed to get slack messages/n")
	}

	// estimate language from original text
	awsClient := NewAwsClient()
	sourceCode, err := awsClient.detectLanguageCode(msg.text)
	if err != nil {
		log.Println("[Error] failed to get language code: ", err)
		return "", err
	}

	// translate text
	translatedText, err := awsClient.translate(msg.text, sourceCode, targetCode)
	if err != nil {
		log.Println("[Error] failed to translate message: ", err)
		return "", err
	}

	// return translation result to slack
	msg.channel = se.Event.Item.Channel
	msg.reaction = se.Event.Reaction
	msg.translated = translatedText
	msg.source = sourceCode
	msg.target = targetCode
	err = c.postMessage(msg)
	if err != nil {
		log.Println("[Error] failed to post slack message: ", err)
		return "", err
	}

	return "", nil
}

// https://api.slack.com/methods/chat.postMessage
func (c *SlackClient) postMessage(msg *slackMsg) error {
	attachment := slack.Attachment{}
	attachment.Pretext = fmt.Sprintf("_The message is translated in_ :%s: _(%s-%s)_", msg.reaction, msg.source, msg.target)
	attachment.Text = msg.translated
	attachment.Footer = msg.text
	attachment.MarkdownIn = []string{"text", "pretext"}

	params := slack.NewPostMessageParameters()
	params.ThreadTimestamp = msg.ts
	params.AsUser = false
	params.Attachments = []slack.Attachment{attachment}

	_, _, err := c.client.PostMessage(msg.channel, "", params)
	if err != nil {
		log.Println("[Error] failed to post slack messages: ", err)
		return err
	}

	return nil
}

// https://api.slack.com/methods/conversations.replies
func (c *SlackClient) getMessage(id string, ts string) (*slackMsg, error) {
	params := &slack.GetConversationRepliesParameters{}
	params.ChannelID = id
	params.Timestamp = ts
	params.Inclusive = true
	params.Limit = 1

	// get slack messages
	msg, _, _, err := c.client.GetConversationReplies(params)
	if err != nil {
		log.Println("[Error] failed to get slack messages: ", err)
		return nil, err
	}

	// get message text
	slMsg := &slackMsg{}
	for _, i := range msg {
		slMsg.ts = i.Timestamp
		if slMsg.ts == "" {
			slMsg.ts = i.ThreadTimestamp
		}

		slMsg.text = i.Text
		if slMsg.text == "" {
			for _, j := range i.Attachments {
				slMsg.text = j.Text
				break
			}
		}
		break
	}

	return slMsg, nil
}
