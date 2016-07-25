package main

// Define the plugin yaml structure

type PluginRuntime struct {
    Service_type string
}

type FlagsOptions
    Help string
    Required bool
}

type PluginDef struct {
    Help string
    Flags map[string]FlagsOptions
}

type PluginAction struct {
    Common PluginDef
    Create PluginDef
    Update PluginDef
    Maintain PluginDef
}

type Plugin struct {
    Runtime PluginRuntime
    Flags PluginAction
}
