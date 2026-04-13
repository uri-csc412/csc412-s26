package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Go version of the config structs.

type LogsConfig struct {
	Directory string `yaml:"dir"`
	Actions   string `yaml:"actions"`
}

type SimulationConfig struct {
	DefaultRows              int  `yaml:"default_rows"`
	DefaultCols              int  `yaml:"default_cols"`
	MinRows                  int  `yaml:"min_rows"`
	MinCols                  int  `yaml:"min_cols"`
	MinThreads               int  `yaml:"min_threads"`
	DefaultThreads           int  `yaml:"default_threads"`
	DrawColoredTravelerHeads bool `yaml:"draw_colored_traveler_heads"`
}

type InkConfig struct {
	MaxLevel     int `yaml:"max_level"`
	MaxAddInk    int `yaml:"max_add_ink"`
	RefillInk    int `yaml:"refill_ink"`
	InitialRed   int `yaml:"initial_red_level"`
	InitialGreen int `yaml:"initial_green_level"`
	InitialBlue  int `yaml:"initial_blue_level"`
}

type TimingConfig struct {
	MinProducerSleepUS  int `yaml:"min_producer_sleep_us"`
	ProducerSleepUS     int `yaml:"producer_sleep_us"`
	InklingSleepUS      int `yaml:"inkling_sleep_us"`
	InklingSleepStepUS  int `yaml:"inkling_sleep_step_us"`
	InklingSleepMinUS   int `yaml:"inkling_sleep_min_us"`
	InklingSleepMaxUS   int `yaml:"inkling_sleep_max_us"`
	ArgWarningSleepS    int `yaml:"arg_warning_sleep_s"`
	ShutdownDelayS      int `yaml:"shutdown_delay_s"`
	EventLoopIntervalMS int `yaml:"event_loop_interval_ms"`
	EventLoopSleepS     int `yaml:"event_loop_sleep_s"`
}

type EventLoopConfig struct {
	MaxRecursiveDepth int `yaml:"max_recursive_depth"`
}

type ProducersConfig struct {
	SpeedupFactorNum  int `yaml:"speedup_factor_num"`
	SpeedupFactorDen  int `yaml:"speedup_factor_den"`
	SlowdownFactorNum int `yaml:"slowdown_factor_num"`
	SlowdownFactorDen int `yaml:"slowdown_factor_den"`
}

type IOConfig struct {
	PipePath string `yaml:"pipe_path"`
}

type IconDirectionsConfig struct {
	North string `yaml:"north"`
	West  string `yaml:"west"`
	South string `yaml:"south"`
	East  string `yaml:"east"`
	None  string `yaml:"none"`
}

type MessagesConfig struct {
	InkTankHeader        string `yaml:"ink_tank_header"`
	LiveThreadsLabel     string `yaml:"live_threads_label"`
	NoArgsMessage        string `yaml:"no_args_message"`
	QuitNoThreadsMessage string `yaml:"quit_no_threads_message"`
	QuitEscMessage       string `yaml:"quit_esc_message"`
	QuitGenericMessage   string `yaml:"quit_generic_message"`
}

type UIConfig struct {
	IconDirections IconDirectionsConfig `yaml:"icon_directions"`
	Messages       MessagesConfig       `yaml:"messages"`
}

type Config struct {
	Logs       LogsConfig       `yaml:"logs"`
	Simulation SimulationConfig `yaml:"simulation"`
	Ink        InkConfig        `yaml:"ink"`
	Timing     TimingConfig     `yaml:"timing"`
	EventLoop  EventLoopConfig  `yaml:"event_loop"`
	Producers  ProducersConfig  `yaml:"producers"`
	IO         IOConfig         `yaml:"io"`
	UI         UIConfig         `yaml:"ui"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
