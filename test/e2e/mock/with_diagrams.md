# Test Document with Mermaid Diagrams

This is a test markdown document that contains multiple Mermaid diagrams.

## Flow Chart

```mermaid
graph TD
    A[Start] --> B[Process]
    B --> C[End]
```

Some text between the diagrams.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant Alice
    participant Bob
    Alice->>Bob: Hello Bob, how are you?
    Bob-->>Alice: I'm good, thanks!
    Alice->>Bob: Goodbye!
```

## Class Diagram

```mermaid
classDiagram
    class Animal {
        +String name
        +int age
        +makeSound()
    }
    class Dog {
        +String breed
        +bark()
    }
    Animal <|-- Dog
```

End of document.
