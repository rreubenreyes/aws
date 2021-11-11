# secret santa reminder app

## what

This is a serverless application which reads (scrapes) a Secret Santa drawpool from [DrawNames](https://www.drawnames.com/), maps people in that pool to a known list of draw participants plus their Discord IDs, and bugs them on Discord until they sign up for the Secret Santa draw.

* Lambda (Go) for application logic
* DynamoDB to lookup participant info
* EventBridge for scheduling

## why though

My friends and I do a Secret Santa drawing every year. Sometimes they forget to sign up for the drawing even though we set it up like 2 months in advance.

