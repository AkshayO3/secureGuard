# SecureGuard

**SecureGuard** is a Go-based security management platform for tracking assets, vulnerabilities, and incidents. It provides robust role-based access control, secure data handling, and a RESTful API built with Gin.

## Features

### Role-Based Access Control (RBAC)
- **Admin**: Full access to all endpoints and operations
- **Analyst**: Create/manage incidents, vulnerabilities, and assets
- **Viewer**: Read-only access to resources

### Management Modules
| Module          | Operations                          |
|-----------------|-------------------------------------|
| **Assets**      | CRUD operations + vulnerability/incident mapping |
| **Vulnerabilities** | CRUD + asset association          |
| **Incidents**   | CRUD + asset association            |

### Data Protection
- **Password Security**:  
  Strong hashing algorithms (bcrypt)
- **Input Validation**:  
  Prevents SQL injection and malformed data
- **Database Security**:  
  Prepared statements for all queries
