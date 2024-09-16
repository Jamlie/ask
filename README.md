# Ask

A CLI tool that gives access to LLMs in the terminal

## Available LLMs

- Gemini 1.5 Flash

## How To Use

- Ask one question:

```
ask gemini "10 reasons why 69 is the best number"
```

- Chat:

```
ask gemini --chat
> what is 0x45 in decimals?
```

## Get Started

1. Install `ask`

```
go install github.com/Jamlie/ask
```

2. Make a `.ask.toml` file in `$HOME` for Windows, macOS and Linux (I'd like to interject for a moment...) in this format:

```toml
gemini_api = "api key"
```
