package flowhash

type FlowHash struct {
	RingCount int
	Key       []byte
	Funcs     map[string]bool
	Table     IndirectTable
}

type Config struct {
	Context RSSContext
}

func NewConfig(opts []Option) (c *Config) {
	c = new(Config)
	for _, opt := range opts {
		opt(c)
	}
	return
}

type Option func(*Config)

func WithRSSContext(context RSSContext) Option {
	return func(c *Config) {
		c.Context = context
	}
}

type SetConfig struct {
	Key     []byte // Sets RSS hash key of the specified network device.
	Func    string // Sets RSS hash function of the specified network device.
	Context RSSContext
	Action  Action
}

func NewSetConfig(opts []SetOption) (c *SetConfig) {
	c = new(SetConfig)

	for _, opt := range opts {
		opt(c)
	}
	return
}

type SetOption func(*SetConfig)

func WithHashKey(key []byte) SetOption {
	return func(c *SetConfig) {
		c.Key = key
	}
}

func WithHashFunc(fn string) SetOption {
	return func(c *SetConfig) {
		c.Func = fn
	}
}

func WithContext(context RSSContext) SetOption {
	return func(c *SetConfig) {
		c.Context = context
	}
}

func NewContext(context RSSContext) SetOption {
	return func(c *SetConfig) {
		c.Context = ETH_RXFH_CONTEXT_ALLOC
	}
}

func WithAction(action Action) SetOption {
	return func(c *SetConfig) {
		c.Action = action
	}
}

func DeleteContext(context RSSContext) SetOption {
	return func(c *SetConfig) {
		c.Context = context
		c.Action = new(Delete)
	}
}
