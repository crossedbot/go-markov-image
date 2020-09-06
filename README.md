## Go Markov Image Generator
Translates an existing image into a new randomly generated image using the
original's colors and color transitions via Markov chaining.
<br/><br/>

### Build
`$ go build -o markov-image`
<br/><br/>

### Run
`$ markov-image input.png output.png -- translate input.png into new image`
<br/><br/>

### API
```
    // Markov is an interface to a markov model of an image
    type Markov interface {
        // GetNextColor returns a randomly selected tranistion color for the
        // given color.
        GetNextColor(c color.Color) color.Color

        // GetRandomColor returns a randomly selected color from the model.
        GetRandomColor() color.Color

        // AddColorTransition adds a color transition, c2, for the given color,
        // c1, to the model.
        AddColorTransition(c1 color.Color, c2 color.Color)

        // Generate returns a new image generated from the current model.
        Generate() *image.RGBA

        // ReadFile reads the given image file and sets the model accordingly.
        ReadFile(filename string) error

        // WriteFile generates a new image based on the current model and writes
        // it to the given file path.
        WriteFile(filename string) error
    }
```