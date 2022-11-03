# Review digest creator

This service polls an iOS app's App Store Connect RSS feed to fetch the app's most recent App Store reviews.

## Usage

Compile and run with `go run .` or use the compiled binary, `digest`. This document will show commands in reference to the compiled binary.

To start the service:
```
./digest
```

The service will run until exited, waiting until a configurable time to publish digests.

Published digests can be found in the `output` folder in a subfolder for the app ID. Each digest is published in json and markdown format. The json file is published for future applications, allowing clients to present their own UI.

## Configuration

The service uses a `config.json` file to configure:

- The time of day at which to publish a digest (default: midnight)
- How many days of reviews to include in a digest (default: 1)
- The app ID of the iOS app to fetch reviews for (default: `595068606`)

By copying this service to different directories and configuring the app ID appropriately, you can use this for as many iOS apps as you wish.

Configure the app using the following command line flags:
- **hour**: Change the hour at which to schedule the digest between 0 and 23.
- **minute**: Change the minute at which to schedule the digest between 0 and 59.
- **interval**: Change the interval (days of reviews) of the digest. Specified in days.
- **app**: Change the app ID for the digest.

Changing the interval is useful for testing purposes.

This command will start the service, configuring the time to 12:30 each day, and setting the app ID to the Uber app:
```
./digest -hour 12 -minute 30 -app 368677368
```

## Generate a digest immediately

Rather than running the service in the background waiting for the next scheduled publication time, you can run the service instantly with the `-now` flag.

Paired with changing the interval, this is great for testing purposes.

```
./digest -now -interval 14
```

Note that if you run this twice in a row, you won't get a second digest. The service will only generate a digest if one does not yet exist in the last 24 hours. Deleting the whole `output` folder is a convenient way to generate new digests for testing purposes. 

## Unit tests

Run tests with `go test`. Due to time constraints, unit tests are sparse, but test the configuration manager and the markdown generation.

## Further considerations

Due to time constraints, pagination is lacking from this implementation. This will work fine for an interval of 1 day, but fetching much longer intervals could miss reviews.
