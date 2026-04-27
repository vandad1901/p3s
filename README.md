# Blog As A Service
Blog As A Service (BAAS) is a simple blogging platform built with Go and Angular. It provides users the ability to spin up their own blog instances with ease, allowing them to focus on creating content rather than managing infrastructure.

## Features
- User authentication and authorization
- Rich text editor for creating and editing blog posts
- Responsive design for mobile and desktop
- RESTful API for managing blog posts and user accounts
- Dockerized for easy deployment

## Architecture
This system uses protobuf as the single source of truth for all service-level contracts. Each service exposes explicitly defined protobuf message types as its outputs, including multiple “view” types when different representations of the same domain data are needed. The REST API and internal service calls both act as consumers of these shared contracts, ensuring consistency across all system boundaries while keeping domain models strictly internal. We effectively elevate contracts to first-class citizens in our architecture. 

## Getting Started
To get started with Blog As A Service, follow these steps:
Open [BAAS](https://baas.vandaddelavari.ir/) in your browser and sign up for an account. Once you have an account, you can assign it a custom domain and start creating your blog posts. You can also explore the API documentation to see how to interact with the platform programmatically.

## Contributing
Contributions to Blog As A Service are welcome! If you have an idea for a new feature or have found a bug, please open an issue or submit a pull request on GitHub.

## License
Blog As A Service is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.