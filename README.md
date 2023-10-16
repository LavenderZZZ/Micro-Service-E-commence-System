# Micro Service E-commerce System (Micro-mall)

## Overview
Our integrated micro-services e-commerce system is built on a range of technologies including Spring Cloud Alibaba, Spring Boot, Oauth, MyBatis, Docker, Jenkins, Kubernetes, Elasticsearch and more. The project business architecture of this system is divided into a front-end mall and a back-end management system, aiming  at provide a complete e-commerce solution. The front-end mall contains four core modules, which are home portal, the product interface, the user interface and the order process, while the back-end management system also contains four key modules, which are the product module, the order module, the marketing module and privilege module, to support comprehensive business requirements.

The Micro-mall System will provide several key benefits, including:
Micro-service Architecture: Adopting micro-service architecture makes the system easier to expand and maintain, and each microservice can be developed, deployed and upgraded independently.
Advanced Technology Stack: Utilizing Spring Cloud Alibaba, Spring Boot  and other leading technologies, the system has high performance, high security and scalability.
User-friendly front-end mall: The front-end mall module provides a powerful user interface, including personalized homepage, product classification, the shopping cart and other functions to enhance user experience.
Efficient Backend Management: The backend management system module supports product, order, marketing and permission management, enabling administrators to efficiently manage business processes.
Search and Analyze: Integrating Elasticsearch, the system has high-speed search and data analysis capabilities to help users quickly find the goods they need and provide data insights.
Containerized Deployment: Adopting Docker technology ensures the consistency and portability of the system in different environments, simplifying the deployment process.
Strong Security: Oauth integration provides advanced authentication and authorization mechanisms to safeguard user data and transactions.
Customization and Flexibility: MyBatis supports customization of database operations to adapt to different business needs.

In conclusion, our integrated micro-services e-commerce system provides a powerful, highly customizable and secure e-commerce platform for businesses. Its adoption of a modern technology stack and flexible architecture gives it multiple advantages that help meet changing market demands and provide a superior experience for customers and administrators.

4. Functional Requirements
4.1 Product Management
Ability to add, update, or remove products with details such as name, price, description, images, and stock information.
Support for product categorization and tagging.

4.2 Search Functionality
Integration with Elasticsearch for fast and relevant product search results.
Support for filtering and sorting the search results.

4.3 User Access Management
The module focuses on managing user creation and access control within the system. The system administrator can create new roles and assign specific access rights based on the organization's needs, ensuring a secure and flexible approach to user access control.
4.3.1 User Registration and Authentication
· The system will provide a user registration process that requires users to enter their personal information, such as username, email address, and password. 
· The backend system will generate a JSON Web Token (JWT) upon successful login that includes user identification and an expiration time.
· The frontend system will store the JWT in the user's browser and use it to authenticate subsequent requests.
· The system will provide an authentication endpoint that validates the JWT and returns a new one with a renewed expiration time.
4.3.2 User Access Control
· The system provides a role-based access control (RBAC) system that enables the creation of different access areas and different access rights.
· The system allows system administrator to create new roles and assign different access rights to them based on the specific needs of the organization.
· The system allows system administrator to assign roles to users, based on their job responsibilities and the access they require to perform their work.
· The system allows system administrator to withdraw roles from users at any time, when their job responsibilities change or when they no longer require access to certain areas of the system.
4.4 Microservices Infrastructure
Distributed system with separate services for user management, product management, order management, etc.
Each service should be able to run, scale, and fail independently.

5. Non-Functional Requirements
5.1 Performance: 
Login Response Time: With 1,000 concurrently logged-in users, the response time for login should not exceed 3 seconds. This metric ensures that even during peak usage times, the user experience remains smooth and frustration-free.
Concurrency Handling: The system's backend and frontend should be optimized to manage a high number of concurrent users and requests. This includes efficient database indexing, caching mechanisms, and load balancing strategies to distribute incoming requests effectively.
Throughput: As a subset of performance, the system should also handle a high number of transactions or actions per second, ensuring that user operations like adding to cart, searching products, or checking out are seamless and swift.

5.2 System Availability: The system should be available for use 24 hours, 7 days a week, with minimal downtime or interruption. 
Uptime: The system should aim for a 99.9% uptime (known as "Three Nines"), which translates to a downtime of not more than 43.8 minutes per month.
Maintenance Windows: Any planned maintenance or updates that may cause downtime should be scheduled during off-peak hours and communicated to users in advance.
Failover Strategy: Implement redundancy measures, such as standby servers or backup databases, to ensure system availability in case of unexpected failures.

5.3 Scalability: The live system will support up to 10000 active user accounts, with 1000 users simultaneously logged in at peak time. The system should be able to adapt to changing usage patterns and requirements over time, without requiring major architectural changes or disruptions. Microservices architecture, when combined with Kubernetes, should allow easy horizontal scaling. 
User Growth: The architecture should be designed to comfortably support up to 10000 active user accounts. Peak load expectations are 10000 simultaneous logins.
Horizontal Scaling with Kubernetes: As demand increases, the system should scale out (add more instances) rather than scale up (increase the capacity of an existing instance). Kubernetes allows for automatic scaling based on predefined metrics such as CPU usage or memory consumption.
Load Balancing: Implement efficient load balancers to distribute incoming traffic across multiple servers or instances, ensuring that no single server is overwhelmed.
5.4 Continuous Integration & Deployment:
Utilize Docker for containerization to ensure consistent environments.
Implement a CI/CD pipeline for automated testing and deployment.
Docker Containerization: By packaging the application and its dependencies into Docker containers, the system ensures consistent environments across development, staging, and production. This eliminates the "it works on my machine" problem.
Automated Testing: As part of the CI pipeline, implement automated testing (unit, integration, and end-to-end tests) to catch regressions and ensure that new features don't break existing functionality.
Deployment Automation: Implement CD tools to automate the deployment process, ensuring that tested and verified code gets deployed to production with minimal human intervention. Rollback strategies should also be in place in case a deployment introduces unforeseen issues.

6. Effort Estimates 
We have planned for a total of four major iterations corresponding to the completion of each phase of the software development life cycle in Agile methodology approach. The total estimated effort for the entire project would be 7 weeks.
Iteration	Duration	Activities	Deliverables
1	1	Requirements Analysis: Work with stakeholders to clearly define the basic requirements and functionality of the system. This may include user stories, use cases, feature lists, etc.
Architecture Design: Design the overall architecture of the system, defining the boundaries of the microservices and how they will communicate with each other.
Develop Iteration Plan: Create a list of tasks for the first iteration and schedule the work to ensure that it is completed in the allotted time.	Project vision document, requirements document, system architecture design document, task list for the first iteration.
2	2.5	Prototyping: Develop a basic system prototype to validate key functionality and user experience.
User Validation: Deliver the prototype to users or key stakeholders for feedback and revision.
Iteration Plan Revision: Update the iteration plan and task list based on user feedback.	Prototyping systems, user feedback reports, revised task lists.
3	2.5	Microservice Development: Concurrent parallel development of microservices, with each microservice team responsible for its specific functionality.
Unit Testing: Developers conduct unit tests to validate the functionality of the microservices.
Integration Testing: Microservices undergo integration testing to ensure that they work together.
System Testing: End-to-end testing of the entire system, simulating real user interactions and scenarios.
Performance testing: Evaluate the performance and scalability of the system and identify potential performance bottlenecks.	Completed microservice code, passed unit tests, integration test reports, system test reports, performance test reports.
4	1	Deployment Preparation: Prepare the production environment and infrastructure, including database configuration, server settings, and so on.
Deploying Microservices: Deploy microservices to the production environment to ensure their stable operation.
Monitoring and Maintenance: Set up monitoring and alerting systems to identify and resolve potential problems.
User Training: Train end-users and administrators to ensure they are able to use the system properly.	Deployed microservices, monitoring and alerting system configuration, user training documentation.
