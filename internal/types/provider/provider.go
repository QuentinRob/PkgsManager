package provider

type Provider string

const (
    Unset  Provider = ""
    Golang Provider = "go"
    APT    Provider = "apt"
    NPM    Provider = "npm"
)
