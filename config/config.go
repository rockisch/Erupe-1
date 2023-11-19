package _config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Mode int

const (
	S1 Mode = iota + 1
	S15
	S2
	S25
	S3
	S35
	S4
	S5
	S55
	S6
	S7
	S8
	S85
	S9
	S10
	F1
	F2
	F3
	F4
	F5
	G1
	G2
	G3
	G31
	G32
	GG
	G5
	G51
	G52
	G6
	G61
	G7
	G8
	G81
	G9
	G91
	G10
	G101
	Z1
	Z2
	ZZ
)

var versionStrings = []string{"S1.0", "S1.5", "S2.0", "S2.5", "S3.0", "S3.5", "S4.0", "S5.0", "S5.5", "S6.0", "S7.0",
	"S8.0", "S8.5", "S9.0", "S10", "FW.1", "FW.2", "FW.3", "FW.4", "FW.5", "G1", "G2", "G3", "G3.1", "G3.2", "GG", "G5",
	"G5.1", "G5.2", "G6", "G6.1", "G7", "G8", "G8.1", "G9", "G9.1", "G10", "G10.1", "Z1", "Z2", "ZZ"}

func (m Mode) String() string {
	return versionStrings[m]
}

// Config holds the global server-wide config.
type Config struct {
	Host                   string `mapstructure:"Host"`
	BinPath                string `mapstructure:"BinPath"`
	Language               string
	DisableSoftCrash       bool     // Disables the 'Press Return to exit' dialog allowing scripts to reboot the server automatically
	HideLoginNotice        bool     // Hide the Erupe notice on login
	LoginNotices           []string // MHFML string of the login notices displayed
	PatchServerManifest    string   // Manifest patch server override
	PatchServerFile        string   // File patch server override
	ScreenshotAPIURL       string   // Destination for screenshots uploaded to BBS
	DeleteOnSaveCorruption bool     // Attempts to save corrupted data will flag the save for deletion
	ClientMode             string
	RealClientMode         Mode
	QuestCacheExpiry       int // Number of seconds to keep quest data cached
	DevMode                bool

	DevModeOptions  DevModeOptions
	GameplayOptions GameplayOptions
	Discord         Discord
	Commands        []Command
	Courses         []Course
	Database        Database
	Sign            Sign
	SignV2          SignV2
	Channel         Channel
	Entrance        Entrance
}

// DevModeOptions holds various debug/temporary options for use while developing Erupe.
type DevModeOptions struct {
	AutoCreateAccount    bool // Automatically create accounts if they don't exist
	CleanDB              bool // Automatically wipes the DB on server reset.
	MaxLauncherHR        bool // Sets the HR returned in the launcher to HR7 so that you can join non-beginner worlds.
	LogInboundMessages   bool // Log all messages sent to the server
	LogOutboundMessages  bool // Log all messages sent to the clients
	LogMessageData       bool // Log all bytes transferred as a hexdump
	MaxHexdumpLength     int  // Maximum number of bytes printed when logs are enabled
	DivaEvent            int  // Diva Defense event status
	FestaEvent           int  // Hunter's Festa event status
	TournamentEvent      int  // VS Tournament event status
	MezFesEvent          bool // MezFes status
	MezFesAlt            bool // Swaps out Volpakkun for Tokotoko
	DisableTokenCheck    bool // Disables checking login token exists in the DB (security risk!)
	QuestDebugTools      bool // Enable various quest debug logs
	EarthStatusOverride  int32
	EarthIDOverride      int32
	EarthMonsterOverride int32
	SaveDumps            SaveDumpOptions
}

type SaveDumpOptions struct {
	Enabled   bool
	OutputDir string
}

// GameplayOptions has various gameplay modifiers
type GameplayOptions struct {
	FeaturedWeapons                int    // Number of Active Feature weapons to generate daily
	MaximumNP                      int    // Maximum number of NP held by a player
	MaximumRP                      uint16 // Maximum number of RP held by a player
	MaximumFP                      uint32 // Maximum number of FP held by a player
	TreasureHuntExpiry             uint32 // Seconds until a Clan Treasure Hunt will expire
	TreasureHuntPartnyaCooldown    uint32 // Seconds until a Partnya can be assigned to another Clan Treasure Hunt
	DisableLoginBoost              bool   // Disables the Login Boost system
	DisableBoostTime               bool   // Disables the daily NetCafe Boost Time
	BoostTimeDuration              int    // The number of minutes NetCafe Boost Time lasts for
	GuildMealDuration              int    // The number of minutes a Guild Meal can be activated for after cooking
	BonusQuestAllowance            uint32 // Number of Bonus Point Quests to allow daily
	DailyQuestAllowance            uint32 // Number of Daily Quests to allow daily
	MezfesSoloTickets              uint32 // Number of solo tickets given weekly
	MezfesGroupTickets             uint32 // Number of group tickets given weekly
	LowLatencyRaviente             bool   // Toggles low latency mode for Raviente, can be network intensive
	RegularRavienteMaxPlayers      uint8
	ViolentRavienteMaxPlayers      uint8
	BerserkRavienteMaxPlayers      uint8
	ExtremeRavienteMaxPlayers      uint8
	SmallBerserkRavienteMaxPlayers uint8
	GUrgentRate                    float32 // Adjusts the rate of G Urgent quests spawning
	GCPMultiplier                  float32 // Adjusts the multiplier of GCP rewarded for quest completion
	GRPMultiplier                  float32 // Adjusts the multiplier of G Rank Points rewarded for quest completion
	GSRPMultiplier                 float32 // Adjusts the multiplier of G Skill Rank Points rewarded for quest completion
	GZennyMultiplier               float32 // Adjusts the multiplier of G Zenny rewarded for quest completion
	MaterialMultiplier             float32 // Adjusts the multiplier of Monster Materials rewarded for quest completion
	ExtraCarves                    uint16  // Grant n extra chances to carve ALL carcasses
	DisableHunterNavi              bool    // Disables the Hunter Navi
	EnableKaijiEvent               bool    // Enables the Kaiji event in the Rasta Bar
	EnableHiganjimaEvent           bool    // Enables the Higanjima event in the Rasta Bar
	EnableNierEvent                bool    // Enables the Nier event in the Rasta Bar
	DisableRoad                    bool    // Disables the Hunting Road
	SeasonOverride                 bool    // Overrides the Quest Season with the current Mezeporta Season
}

// Discord holds the discord integration config.
type Discord struct {
	Enabled           bool
	BotToken          string
	RealtimeChannelID string
}

// Command is a channelserver chat command
type Command struct {
	Name    string
	Enabled bool
	Prefix  string
}

// Course represents a course within MHF
type Course struct {
	Name    string
	Enabled bool
}

// Database holds the postgres database config.
type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Sign holds the sign server config.
type Sign struct {
	Enabled bool
	Port    int
}

// SignV2 holds the new sign server config
type SignV2 struct {
	Enabled     bool
	Port        int
	PatchServer string
}

type Channel struct {
	Enabled bool
}

// Entrance holds the entrance server config.
type Entrance struct {
	Enabled bool
	Port    uint16
	Entries []EntranceServerInfo
}

// EntranceServerInfo represents an entry in the serverlist.
type EntranceServerInfo struct {
	IP          string
	Type        uint8  // Server type. 0=?, 1=open, 2=cities, 3=newbie, 4=bar
	Season      uint8  // Server activity. 0 = green, 1 = orange, 2 = blue
	Recommended uint8  // Something to do with server recommendation on 0, 3, and 5.
	Name        string // Server name, 66 byte null terminated Shift-JIS(JP) or Big5(TW).
	Description string // Server description
	// 4096(PC, PS3/PS4)?, 8258(PC, PS3/PS4)?, 8192 == nothing?
	// THIS ONLY EXISTS IF Binary8Header.type == "SV2", NOT "SVR"!
	AllowedClientFlags uint32

	Channels []EntranceChannelInfo
}

// EntranceChannelInfo represents an entry in a server's channel list.
type EntranceChannelInfo struct {
	Port           uint16
	MaxPlayers     uint16
	CurrentPlayers uint16
}

var ErupeConfig *Config

func init() {
	var err error
	ErupeConfig, err = LoadConfig()
	if err != nil {
		preventClose(fmt.Sprintf("Failed to load config: %s", err.Error()))
	}
}

// getOutboundIP4 gets the preferred outbound ip4 of this machine
// From https://stackoverflow.com/a/37382208
func getOutboundIP4() net.IP {
	conn, err := net.Dial("udp4", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.To4()
}

// LoadConfig loads the given config toml file.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("DevModeOptions.SaveDumps", SaveDumpOptions{
		Enabled:   true,
		OutputDir: "save-backups",
	})

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	if c.Host == "" {
		c.Host = getOutboundIP4().To4().String()
	}

	for i := range versionStrings {
		if strings.ToUpper(c.ClientMode) == versionStrings[i] {
			c.RealClientMode = Mode(i + 1)
			c.ClientMode = strings.ToUpper(c.ClientMode)
			if c.RealClientMode <= G101 {
				c.ClientMode += " (Debug only)"
			}
		}
	}
	if c.RealClientMode == 0 {
		c.ClientMode = versionStrings[len(versionStrings)-1]
		c.RealClientMode = ZZ
	}

	return c, nil
}

func preventClose(text string) {
	if ErupeConfig.DisableSoftCrash {
		os.Exit(0)
	}
	fmt.Println("\nFailed to start Erupe:\n" + text)
	go wait()
	fmt.Println("\nPress Enter/Return to exit...")
	fmt.Scanln()
	os.Exit(0)
}

func wait() {
	for {
		time.Sleep(time.Millisecond * 100)
	}
}
