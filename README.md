# Retrosheet
!This projet is a work in progress!

I wanted to write a couple of applications to download, and load the data into a database for replaying baseball game plays for various reasons.

## Building Applications and Setup
The applications are written in GO
* downloader
* dbloader (MySQL)
* schema (Creates DB and tables)
<pre>git clone [REPO]</pre>

or 

<pre>go get github.com/wazupwiddat/retrosheet-go</pre>

Install dependencies
<pre>dep ensure </pre>

Build applications
<pre>make all [build_down, build_loader, build_migration]</pre>

MySQL install
<pre>docker run --name baseball -p 3306:3306 -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mysql:5.6</pre>

Applying Schema to MySQL DB
<pre>./bin/retrosheet-migration "root:@tcp(localhost:3306)/baseball" up</pre>

Dropping tables
<pre>./bin/retrosheet-migration "root:@tcp(localhost:3306)/baseball" down</pre>

## Running Applications
Downloader
<pre>./bin/retrosheet-downloader</pre>

This application will simply download all the ZIP files from the retrosheet site skipping all invalid years.  Parameters are listed below and the defaults will download the entire lot.

<pre>Usage of ./bin/retrosheet-downloader:
  -end int
    	Start year. Default: 1921 (default 2019)
  -output string
    	Download output path. Default: '.' (default "output")
  -start int
    	Start year. Default: 1921 (default 1921)
</pre>

DB Loader
<pre>./bin/retrosheet-dbloader -output output</pre>

This application will read the `-output output` directory (or where you downloaded the ZIP files) and load everything into the MySQL DB.

**Note: if you are going to load all the data in you will need ~3G in storage space, a fast'ish computer, and about 4 hours depending on hardware.

## Data Models

`teams`
<pre>type Team struct {
	ID       int
	TeamCode string `db:"team_code"`
	Year     int
	Name     string
	Mascot   string
	League   League
}</pre>
`games`
<pre>type Game struct {
	ID      int
	GameID  string `db:"game_id"`
	Visitor int
	Home    int
	Played  time.Time
}
</pre>
`players`
<pre>type Player struct {
	ID        int
	PlayerID  string `db:"player_id"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	Bats      Handed
	Throws    Handed
}</pre>
`events`
<pre>type GameEvent struct {
	ID         int         `db:"id" json:"-"`
	GameID     int         `db:"game_id" json:"-"`
	Event      EventType   `db:"event" json:"-"`
	Inning     int         `db:"inning" json:"-"`
	InningHalf InningHalf  `db:"inning_half" json:"-"`
	Player     int         `db:"player_id" json:"-"`
	Play       EventDetail `db:"-" json:"event_detail"`
	PlayJSON   string      `db:"event_detail" json:"-"`
}

type EventDetail struct {
	Play       BasicPlay
	ExtraPlays []BasicPlay
	Fielders   []Position
	Modifiers  []Modifier
	RunnerAdv  []RunnerAdvance
}</pre>

## Notes
* after loading the data into the database, it would be helpful to add a few indexes
> 
<pre>
CREATE TABLE `game_events` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_id` int(11) NOT NULL,
  `player_id` int(11) NOT NULL,
  `inning` int(11) NOT NULL,
  `inning_half` int(11) NOT NULL,
  `event` int(11) NOT NULL,
  `event_detail` text,
  PRIMARY KEY (`id`),
  KEY `game_id` (`game_id`),
  KEY `inning` (`inning`,`inning_half`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
</pre>
>  KEY `game_id` (`game_id`)

>  KEY `inning` (`inning`,`inning_half`)

* `event_detail` - column in `game_events` is a JSON document
<pre>
{
  "event_detail": {
    "Play": 12,
    "ExtraPlays": [],
    "Fielders": [
      4
    ],
    "Modifiers": [],
    "RunnerAdv": [
      {
        "StartBase": 2,
        "FinishBase": 3
      }
    ]
  }
}</pre>