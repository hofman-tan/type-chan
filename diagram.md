```mermaid
---
title: Class diagram
---
classDiagram
    direction BT
    class App {
      currentPage
      changePage(Page)
    }

    class Page {
      <<interface>>

      Init() tea.Cmd
      Update(tea.Msg) tea.Cmd
      View() string
    }
    class TypingPage {
      app
      currentState
      changeState(State)
    }
    class Text {
    }
    class TypingPageViewBuilder {
      render()
    }
    class ResultPage{
      app
    }
    
    class QuoteFetcher {
    }

    class Timer {
      <<interface>>
    }
    class CountUpTimer {
    }
    class CountDownTimer {
    }

    class State {
      <<interface>> 

      handleLetter(string)
      handleSpace()
      handleBackspace()
    }
    class CorrectState {
      typingPage
    }
    class WrongState {
      typingPage
    }

    App --> Page
    TypingPage --|> Page
    ResultPage --|> Page
    TypingPage <..> ResultPage : transitions to
    TypingPage --> TypingPageViewBuilder

    TypingPage --> State
    CorrectState --|> State
    WrongState --|> State
    CorrectState <..> WrongState : transitions to

    TypingPage --> Text
    TypingPage ..> QuoteFetcher
    TypingPage --> Timer
    CountUpTimer --|> Timer
    CountDownTimer --|> Timer
```