# Blog Platform
A simple blogging platform built with Go and Angular. It allows users to create and manage blog content with support for media uploads and external integrations, focusing on a clean service-oriented backend architecture.

## Features
User authentication and authorization via dedicated auth service (JWT + refresh tokens)
Rich text editor for creating and editing blog posts
Media uploads (images, video, audio) with asynchronous processing
Cross-posting integration with external platforms (e.g. Telegram channels)
Responsive design for mobile and desktop
RESTful API for managing blog posts, media, and user accounts
Dockerized for easy deployment

## Architecture
This system uses protobuf as the single source of truth for all service-level contracts. Each service exposes explicitly defined protobuf message types as its outputs, including multiple “view” types when different representations of the same domain data are needed. The REST API and internal service calls both act as consumers of these shared contracts, ensuring consistency across all system boundaries while keeping domain models strictly internal. We effectively elevate contracts to first-class citizens in our architecture.

In addition to contract-driven service boundaries, the system is designed as a distributed, event-driven architecture:

A dedicated Auth Service is responsible for identity management, authentication, and session control using JWT access tokens and rotating refresh tokens.
Core application services (Blog, Media, Integrations) operate as independent services and validate access tokens locally without coupling to the auth service at request time.
An event bus (Kafka/RabbitMQ) is used for asynchronous workflows such as media processing, cross-posting, and downstream integrations.
Media processing is handled by worker services consuming events from the queue, enabling scalable and fault-tolerant processing of uploads.
External integrations (e.g. Telegram) are implemented as isolated consumer services reacting to domain events.

This separation of concerns enables horizontal scalability, independent deployment of services, and clear ownership boundaries across the system.

## Getting Started
To get started with the Blog Platform, follow these steps:
1. Clone the repository:
   ```bash
   git clone
    ```
2. Navigate to the project directory:
   ```bash
   cd blog-platform
   ```
3. run the setup script to initialize the environment and dependencies:
   ```bash
   just dev-setup
   ```
4. Build and run the application using Docker Compose:
   ```bash
   just compose-up
   ``` 
4. Access the application at `http://localhost:4200` and start creating your blog posts!

## Contributing
Contributions to Blog As A Service are welcome! If you have an idea for a new feature or have found a bug, please open an issue or submit a pull request on GitHub.

## License
Blog As A Service is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.