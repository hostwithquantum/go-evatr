## Contributing

All contributions (bug fixes, features, docs) are welcome! Please feel free to submit a Pull Request.

Ensure to run the tests for all contributions, please try to keep backwards compatibility.

```bash
go vet ./...
go test -v ./...
```

In addition, all features and bug fixes should be accompanied by a new test case to demonstrate how the code was not working before and how the new code improves it.