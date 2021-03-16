# ElectricKiwiApi
Unofficial API for Electric Kiwi written in Go.
This is not supported by Electric Kiwi at all and may stop working at any point.

### Example usage
```go
sess, err := ElectricKiwiApi.NewSession("email", "password")
if err != nil {
    log.Fatal(err)
}

h, err := ElectricKiwiApi.NewHourOfPower("2100") // only use supported hours from Electric Kiwi in 24 hour format e.g. "1130", "1200", "1230", "1300", "1330", ...
if err != nil {
    return err
}

err = sess.UpdateHourOfPower(h)
if err != nil {
    return err
}

return nil
```

### Supported Features
The API only supports setting the hour of power. This is useful considering that
Electric Kiwi don't officially support automatically setting the hour of power on
a per day basis. You can wrap it in a schedule. Honestly I don't know what else you'd
want the API to do.

### Contributing
Please feel free to make a PR against this repo or make a feature request.