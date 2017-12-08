# Customerio

A golang client for the [Customer.io](https://customer.io) [REST API](https://learn.customer.io/api).
*Tested with Go1.4*

Godoc here: [https://godoc.org/github.com/customerio/go-customerio](https://godoc.org/github.com/customerio/go-customerio)


## Installation

Add this line to your application's imports:

```go
import (
    // ...
    "github.com/customerio/go-customerio"
)
```

And then execute:

    go get

Or install it yourself:

    $ go get "github.com/customerio/go-customerio"

## Usage

### Before we get started: API client vs. JavaScript snippet

It's helpful to know that everything below can also be accomplished
through the [Customer.io JavaScript snippet](https://learn.customer.io/developer-documentation/javascript-quick-start.html).

In many cases, using the JavaScript snippet will be easier to integrate with
your app, but there are several reasons why using the API client is useful:

* You're not planning on triggering emails based on how customers interact with
  your website (e.g. users who haven't visited the site in X days)
* You're using the javascript snippet, but have a few events you'd like to
  send from your backend system.  They will work well together!
* You'd rather not have another javascript snippet slowing down your frontend.
  Our snippet is asynchronous (doesn't affect initial page load) and very small, but we understand.

In the end, the decision on whether or not to use the API client or
the JavaScript snippet should be based on what works best for you.
You'll be able to integrate **fully** with [Customer.io](https://customer.io) with either approach.

### Setup

Create an instance of the client with your [Customer.io](https://customer.io) credentials
which can be found on the [Customer.io integration screen](https://fly.customer.io/account/customerio_integration).

```go
cio := customerio.NewCustomerIO("YOUR SITE ID", "YOUR API SECRET KEY")
```

### Identify logged in customers

Tracking data of logged in customers is a key part of [Customer.io](https://customer.io). In order to
send triggered emails, we must know the email address of the customer.  You can
also specify any number of customer attributes which help tailor [Customer.io](https://customer.io) to your
business.

Attributes you specify are useful in several ways:

* As customer variables in your triggered emails.  For instance, if you specify
the customer's name, you can personalize the triggered email by using it in the
subject or body.

* As a way to filter who should receive a triggered email.  For instance,
if you pass along the current subscription plan (free / basic / premium) for your customers, you can
set up triggers which are only sent to customers who have subscribed to a
particular plan (e.g. "premium").

You'll want to indentify your customers when they sign up for your app and any time their
key information changes. This keeps [Customer.io](https://customer.io) up to date with your customer information.

```go
// Arguments
// customerID (required) - a unique identifier string for this customers
// attributes (required) - a ```map[string]interface{}``` of information about the customer. You can pass any
//                         information that would be useful in your triggers. You 
//                         should at least pass in an email, and created_at timestamp.
//                         your interface{} should be parseable as Json by 'encoding/json'.Marshal

cio.Identify("5", map[string]interface{}{
  "email": "bob@example.com",
  "created_at": time.Now().Unix(),
  "first_name": "Bob",
  "plan": "basic",
})
```

### Deleting customers

Deleting a customer will remove them, and all their information from
Customer.io.  Note: if you're still sending data to Customer.io via
other means (such as the javascript snippet), the customer could be
recreated.

```go
// Arguments
// customerID (required) - a unique identifier for the customer.  This
//                         should be the same id you'd pass into the
//                         `identify` command above.

cio.Delete("5")
```

### Tracking a custom event

Now that you're identifying your customers with [Customer.io](https://customer.io), you can now send events like
"purchased" or "watchedIntroVideo".  These allow you to more specifically target your users
with automated emails, and track conversions when you're sending automated emails to
encourage your customers to perform an action.


```go
// Arguments
// customerID (required)  - the id of the customer who you want to associate with the event.
// name (required)        - the name of the event you want to track.
// attributes (optional)  - any related information you'd like to attach to this
//                          event, as a ```map[string]interface{}```. These attributes
//                          can be used in your triggers to control who should receive
//                          the triggered email. You can set any number of data values.

cio.Track("5", "purchase", map[string]interface{}{
    "type": "socks",
    "price": "13.99",
})
```

### Trigger a broadcast campaign

The trigger method is used to trigger API Triggered Broadcast campaigns in Customer.io. When sending data along with your trigger, it is required to send a key/value pair in the data map.

```go
// Arguments
// campaignID (required)  - the id of the broadcast campaign you want to trigger.
// attributes (optional)  - any related information you'd like to attach to this
//                          trigger, as a ```map[string]interface{}```. These attributes
//                          can be used in your message templates. You can set any number
//                          of data values.
// recipients (optional)  - override the recipient setting for your campaign.

cio.Trigger("5", map[string]interface{}{
    "title": "Roadrunner spotted in Albuquerque!",
    "date": 1511315635,
    "text": "We received reports of a roadrunner in your immediate area! Head to your dashboard to view more information!",
    "recipients": map[string]interface{}{
        "segment": map[string]interface{}{
            "id": 7,
        },
    },
})
```

## Contributing

1. Fork it
2. Clone your fork (`git clone git@github.com:MY_USERNAME/go-customerio.git && cd go-customerio`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes (`git commit -am 'Added some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request
