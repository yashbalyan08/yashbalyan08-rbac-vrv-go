# **Role-Based Access Control (RBAC) System**

## **Project Overview**

This project demonstrates the implementation of **Authentication**, **Authorization**, and **Role-Based Access Control (RBAC)** using Go as the backend, PostgreSQL as the database, and a basic frontend created with HTML and JavaScript. The system ensures secure user management and controlled access to resources based on user roles.

---

## **Features**

1. **User Authentication**
   - Secure user registration, login, and logout.
   - Passwords hashed using industry-standard encryption.
   - JWT-based session management.

2. **Authorization with RBAC**
   - Role-based access to endpoints and resources.
   - Defined roles: Admin, User, Moderator.
   - Permissions mapped to roles for fine-grained control.

3. **Frontend**
   - Basic HTML and JavaScript interface for user interaction.
   - Supports user login and access to role-specific resources.

4. **Security Best Practices**
   - Implements secure password hashing and JWT token storage.
   - Prevents unauthorized access with middleware checks for roles.

---

## **Technologies Used**

- **Backend:** Go  
- **Database:** PostgreSQL  
- **Frontend:** HTML, JavaScript  
- **Authentication:** JWT  
- **Database Management:** GORM (Go ORM for PostgreSQL)

---

## **System Architecture**

1. **Models**
   - **User:** Stores user details and hashed passwords.
   

2. **Endpoints**
   - `/register`: User registration.
   - `/login`: User login and token generation.
   - `/logout`: User logout and token invalidation.
   - `/dashboard`: Only accessible to Admin.
   - `/profile`: Only accessible to User.

3. **Middleware**
   - Auth middleware to validate JWT tokens.
   - RBAC middleware to check role-based permissions.

---

## **Setup and Installation**

### **Prerequisites**  
Before running the project, ensure you have the following installed on your system:  

- [Go](https://go.dev/doc/install) (version 1.23.1 or above recommended)  
- [PostgreSQL](https://www.postgresql.org/download/) (version 12 or above)  


---

## **Steps to Run Locally**

1. **Clone the Repository**  
   Clone this repository to your local machine:  
   ```bash
   git clone https://github.com/yashbalyan08/yashbalyan08-rbac-vrv-go.git
   cd yashbalyan08-rbac-vrv-go
2. **Set Up the Database**  
   - Directly use the docker-compose file to run the postgres image
   ```bash
    $ docker-compose up -d
   ```  
   - Update the database connection details in `config/db.go` to match your environment:  
     ```go
     package config

     import (
         "gorm.io/driver/postgres"
         "gorm.io/gorm"
         "log"
     )
     
     var DB *gorm.DB

     const (
         Host     = os.Getenv("DB_HOST")
         Port     = os.Getenv("DB_PORT")
         User     = os.Getenv("DB_USER")
         Password = os.Getenv("DB_PASSWORD")
         DbName   = os.Getenv("DB_NAME")

     )

     func InitDB() {
         dsn := "host=" + Host + " user=" + User + " password=" + Password + " dbname=" + DbName + " port=" + string(Port) + " sslmode=disable TimeZone=UTC"
         DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
         if err != nil {
             log.Fatal("Failed to connect to the database:", err)
         }
     }
     ```

   - Ensure the `config/ConnectDB()` function is invoked in your main entry point (`main.go`) to establish the database connection.

3. **Hashing To Portect Password**  
   - Hashing the password using `bcrypt` while saving
    Example code:  
     ```go
        package models

        import "golang.org/x/crypto/bcrypt"

        type User struct {
	        ID       uint   `gorm:"primaryKey"`
	        Username string `gorm:"unique"`
	        Password string
	        Role     string
        }

        func HashPassword(password string) (string, error) {
	        bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	        return string(bytes), err
        }

        
     ```

   - Comapring the password to the saved hashed password
     ```go
     func CheckPassword(hashedPass, password string) bool {
	        err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	        return err == nil
        }
     ```

4. **Start the Backend**  
   - Run the application with the following command:  
     ```makefile
     make
     make run
     ```
   - The server should start on `http://localhost:8080`.

---

# API Routes

### **Authentication and Registration**
- **GET** `/login`  
  - Description: Displays the login page.  
  - Middleware: None.  
  - Response: Renders `login.html`.  

- **GET** `/register`  
  - Description: Displays the registration page.  
  - Middleware: `AuthMiddleware`.  
  - Response: Renders `register.html`.  

- **POST** `/register`  
  - Description: Handles user registration.  
  - Middleware: None.  
  - Controller: `controllers.Register`.  

- **POST** `/login`  
  - Description: Handles user login.  
  - Middleware: None.  
  - Controller: `controllers.Login`.  

### **Authenticated Routes**
- **GET** `/profile`  
  - Description: Displays user profile data (role: User).  
  - Middleware:  
    - `AuthMiddleware`: Ensures the user is authenticated.  
    - `AuthorizeRoles`: Restricts access to users with the "User" role.  
  - Response: JSON object with user profile data.  

  ```json
  {
      "message": "User profile data"
  }
  ```

- **GET** `/dashboard`  
  - Description: Displays the admin dashboard (role: Admin).  
  - Middleware:  
    - `AuthMiddleware`: Ensures the user is authenticated.  
    - `AuthorizeRoles`: Restricts access to users with the "Admin" role.  
  - Response: JSON object with dashboard data.  

  ```json
  {
      "message": "Dashboard for authenticated users"
  }
  ```

- **GET** `/logged-in`  
  - Description: Displays a logged-in status page.  
  - Middleware: 
    - `AuthMiddleware`:Ensures the user is authenticated.
    - `AuthorizeRoles`: Restricts access to users with the "User" and "Admin" role.
  - Response: Renders `logged-in.html`.  

### **Logout**
- **POST** `/logout`  
  - Description: Handles user logout.  
  - Middleware: None.  
  - Controller: `controllers.Logout`.  

---

## Usage

### 1. Register a User  
Navigate to the registration page (`/register`) to create a new user account. Provide the required details, such as username, password, and any additional fields.

### 2. Login  
Access the login page (`/login`) and enter your credentials. Upon successful authentication, a JWT token will be issued, which must be used to access protected resources.

### 3. Access Resources  
Use the issued JWT token to access role-specific endpoints:
- **User Role:** Access `/profile` to view user-specific information.  
- **Admin Role:** Access `/dashboard` to view admin-specific data.

Unauthorized access attempts will result in a **403 Forbidden** response. Ensure the token is included in the `Authorization` header of your API requests.

### 4. Logout  
Send a POST request to `/logout` to invalidate the session and log out securely.
---

## Author

This project was developed by **Yash Balyan**.  
For any inquiries or contributions, feel free to reach out via [LinkedIn](https://linkedin.com/in/yash-balyan) or visit my [GitHub](https://github.com/yashbalyan08).
