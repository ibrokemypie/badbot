# badbot
## a bot for discord written in golang


### Usage
1. Download source
2. Install golang
3. Build binary  wuth ``go build``
4. [Configure]
5. Run with ``./badbot``


### Configure
Configuration is done with a toml file called ``config.toml`` with the following contents:
```
token = "the bot's token to authenticate" string
trustedusers = "array of user ids to allow trusted commands to" []string
nickname = "the bot's nickname" string
image = "base64 encoded image for profile picture" string
status = "the bot's status" string

replies = "whether to enable replies plugin" bool

engineid = "engineid of google custom search engine" string
googleapi = "your google custom search engine api key" string
youtubeapi = "your google api key with youtube search enabled" string
lastfmapi = "your lastfm api key" string
lastfmuser = "your lastfm username" string

```


### Commands
``>ping`` replies with pong

``>pong`` replies with ping

``>sarahah`` replies with my sarahah

``>qoohme`` replies with my qoohme

``>git`` replies with this repo

``>help`` gives help



``>woof`` sends a dog from random.dog

``>meow`` sends a cat from random.cat


``>playing`` sends the currently playing track on last.fm



``>>> name quote`` adds a quote under that name

``>> name`` sends a random quote under that name

``>qid id`` sends the quote with given id

``>qdel id`` removes quote with given id


#### trusted user only commands
``>search term`` searches google for term and displays 10 pages of results

``>yt term`` searches youtube for term and displays 10 pages of results


``>say words`` makes the bot say words

``>game words`` sets the bot's game to words

``>status words`` sets the bot's status to words

``>pfp <attached image>`` converts the attached image to base64 and sets it as the profile picture
