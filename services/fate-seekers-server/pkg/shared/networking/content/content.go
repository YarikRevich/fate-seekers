package content

// NetworkingContentChannel represent networking content channel interface.
type NetworkingContentChannel interface {
	// Schedule performs channel call once.
	Schedule(args interface{}, finishCallback func())
}

// err := udpt.SendString("127.0.0.1:9876", "main", "Hello World!", cryptoKey)
// if err != nil {
//     fmt.Println("failed sending:", err)
// }
