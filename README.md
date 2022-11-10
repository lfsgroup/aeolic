# aeolic

Opinionated approach to sending slack requests to a channel: 

(using a simple http post request, no slack apis)

1. You provide the BlockBuild json [via block kit builder](https://app.slack.com/block-kit-builder/T02K6GVUGAY#%7B%22blocks%22:%5B%7B%22type%22:%22section%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22Hello,%20Assistant%20to%20the%20Regional%20Manager%20Dwight!%20*Michael%20Scott*%20wants%20to%20know%20where%20you'd%20like%20to%20take%20the%20Paper%20Company%20investors%20to%20dinner%20tonight.%5Cn%5Cn%20*Please%20select%20a%20restaurant:*%22%7D%7D,%7B%22type%22:%22divider%22%7D,%7B%22type%22:%22section%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22*Farmhouse%20Thai%20Cuisine*%5Cn:star::star::star::star:%201528%20reviews%5Cn%20They%20do%20have%20some%20vegan%20options,%20like%20the%20roti%20and%20curry,%20plus%20they%20have%20a%20ton%20of%20salad%20stuff%20and%20noodles%20can%20be%20ordered%20without%20meat!!%20They%20have%20something%20for%20everyone%20here%22%7D,%22accessory%22:%7B%22type%22:%22image%22,%22image_url%22:%22https://s3-media3.fl.yelpcdn.com/bphoto/c7ed05m9lC2EmA3Aruue7A/o.jpg%22,%22alt_text%22:%22alt%20text%20for%20image%22%7D%7D,%7B%22type%22:%22section%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22*Kin%20Khao*%5Cn:star::star::star::star:%201638%20reviews%5Cn%20The%20sticky%20rice%20also%20goes%20wonderfully%20with%20the%20caramelized%20pork%20belly,%20which%20is%20absolutely%20melt-in-your-mouth%20and%20so%20soft.%22%7D,%22accessory%22:%7B%22type%22:%22image%22,%22image_url%22:%22https://s3-media2.fl.yelpcdn.com/bphoto/korel-1YjNtFtJlMTaC26A/o.jpg%22,%22alt_text%22:%22alt%20text%20for%20image%22%7D%7D,%7B%22type%22:%22section%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22*Ler%20Ros*%5Cn:star::star::star::star:%202082%20reviews%5Cn%20I%20would%20really%20recommend%20the%20%20Yum%20Koh%20Moo%20Yang%20-%20Spicy%20lime%20dressing%20and%20roasted%20quick%20marinated%20pork%20shoulder,%20basil%20leaves,%20chili%20&%20rice%20powder.%22%7D,%22accessory%22:%7B%22type%22:%22image%22,%22image_url%22:%22https://s3-media2.fl.yelpcdn.com/bphoto/DawwNigKJ2ckPeDeDM7jAg/o.jpg%22,%22alt_text%22:%22alt%20text%20for%20image%22%7D%7D,%7B%22type%22:%22divider%22%7D,%7B%22type%22:%22actions%22,%22elements%22:%5B%7B%22type%22:%22button%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Farmhouse%22,%22emoji%22:true%7D,%22value%22:%22click_me_123%22%7D,%7B%22type%22:%22button%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Kin%20Khao%22,%22emoji%22:true%7D,%22value%22:%22click_me_123%22,%22url%22:%22https://google.com%22%7D,%7B%22type%22:%22button%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Ler%20Ros%22,%22emoji%22:true%7D,%22value%22:%22click_me_123%22,%22url%22:%22https://google.com%22%7D%5D%7D%5D%7D), as a golang template
2. You call the interface [method](cmd/slack/main.go), with the auth, channel id and payload required (see examples below)
3. The aeloic code makes the slack call
4. Profit

<br>

----

<br>

## Getting started


<br>

### Adding Templates
Create a directory and add files with the suffix `tmpl.json`. When creating a new instance of aeolic, pass in the directory path, aeolic will automatically load all files matching the suffix.

template names are the names of the file without the suffix, for example:

file: `basic.tmpl.json`
template name: `basic`

<br>

----

<br>


### Making the slack call

Create a new instance with the slack token and the path to your templates directory.


```golang

	c, err := aeolic.New("<your-slack-token>", "<path/to/template/dir>")

	if err != nil {
        // handle error
	}

```

Make the api call

```golang


	if err := c.SendMessage("<your-slack-channel-id>", "<template-name>", map[string]string{
		"hello": "world",
	}); err != nil {
        // handle error
	}

```

<br>

### Provide your own template map

<br>

Using go [embed](https://pkg.go.dev/embed) feature you can provide your own template map, the template name will be the key.

```golang
import (
        _ "embed"
)

//go:embed templates/basic.tmpl.json
var basicTemplate string

func main() {

	customMap := map[string]string{
		"basic": basicTemplate,
	}

	c := aeolic.NewWithMap(<token>, customMap)

	if err := c.SendMessage(channel, "basic", map[string]string{
		"user_name": "Allan Bond",
	}); err != nil {
		log.Fatal("failed ", err)
	}

}
```
<br>

### Provide your own file system
Using go [embed](https://pkg.go.dev/embed) feature you can provide your own file system.

```go
import (
	_ "embed"
)
//go:embed templates/*.tmpl.json
var content embed.FS

func main() {

	c := aeolic.NewWithFS(<token>, content, "templates")

	if err := c.SendMessage(channel, "basic", map[string]string{
		"user_name": "Allan Bond",
	}); err != nil {
		log.Fatal("failed ", err)
	}

}
```

<br>

----

<br>

## Error Handling

<br>


More context is provided if it's an api error

```golang

	if err := c.SendMessage("<your-slack-channel-id>", "<template-name>", map[string]string{
		"hello": "world",
	}); err != nil {
                var apiErr *aeolic.APIError
                if errors.As(err, &apiErr) {
                        // non 2xx,3xx response for example: 
                        // StatusCode: 400
                        // StatusText: Bad Request
                        // Message: "invalid_blocks"
                        // Context: "https://api.slack.com/methods/chat.postMessage#errors"

                
                        /** ... */
                }

                // handle other errors
                fmt.Println(err)
	}

```


