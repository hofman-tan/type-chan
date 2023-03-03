```mermaid
---
title: Type-chan
---
classDiagram
    direction BT
    class App {
      currentPage
      changePage(Page)
    }

    class Page {
      <<interface>>

      app

      Init() tea.Cmd
      Update(tea.Msg) tea.Cmd
      View() string
    }
    class TypingPage {
      app
      currentState
      changeState(State)
    }
    class ResultPage{
      app
    }

    class State {
      <<interface>> 

      typingPage

      handleLetter(string)
      handleSpace()
      handleBackspace()
      view() string
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
    TypingPage <--> ResultPage : changes between

    TypingPage --> State
    CorrectState --|> State
    WrongState --|> State
    CorrectState <--> WrongState : switches between
```