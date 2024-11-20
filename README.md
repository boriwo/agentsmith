# agentsmith

An LLM enabled chatbot acting as a conversational assistant on top of 
domain specific knowledge bases of given facts. Currently, this implementation
is only for experimentation and demonstration purposes. It relies heavily on the
Embed() and ranking to find the best possible answer from the knowledge base
for a given question. There are also system commands to explore and modify
the knowledge bases and these system commands can also be accessed through
the same natural language chat interface. This is all still very much work
in progress.

## demo

You can interact with the bot via Slack, CLI or web interface at 
`http://localhost:8080/agentsmith`.

```
Q: list availabe knowledge bases
```

```
A: startrek, starwars
```

```
Q: select knowledge base startrek
```

```
A: set current knowledge base to startrek
```

```
Q: who are some of the main characters?
```

```
A: The following are notable characters of the original series:Captain James T. Kirk (William Shatner) - The charismatic captain of the USS Enterprise.Mr. Spock (Leonard Nimoy) - A half-human, half-Vulcan science officer known for his logical thinking.Dr. Leonard 'Bones' McCoy (DeForest Kelley) - The ship's chief medical officer, known for his emotional and often cantankerous demeanor.Other notable crew members include Lieutenant Uhura, Mr. Sulu, Ensign Chekov, and Chief Engineer Scott.
```

```
Q: now select starwars as your knowledge base!
```

```
A: set current knowledge base to starwars
```

```
Q: please list some of the characters!
```

```
A: Key characters include Luke Skywalker, Princess Leia, Han Solo, Darth Vader, Yoda, and Obi-Wan Kenobi.
```



