# Tempo: A Productivity App for To-Do Lists, Notes, and Journaling

## 1. Introduction

Tempo is a versatile productivity application designed to help users organize their tasks, capture their thoughts, and reflect on their days. It provides a seamless experience across both a web application and a dedicated Android app, ensuring that your to-do lists, notes, and journal entries are always in sync and accessible. The backend is built with the high-performance Go programming language and utilizes a robust PostgreSQL database to securely store user data.

## 2. Core Features

Tempo is built around three core pillars of personal productivity:

### To-Do Lists
*   **Task Management:** Create, edit, and delete tasks with due dates and priority levels.
*   **Project Organization:** Group tasks into projects or categories for better organization.
*   **Subtasks:** Break down larger tasks into smaller, more manageable subtasks.
*   **Reminders:** Set reminders for important deadlines.
*   **Filtering and Sorting:** Easily filter and sort tasks by project, priority, or due date.

### Notes
*   **Rich Text Editing:** A simple and intuitive editor that supports basic formatting like headers, bold, italics, and lists.
*   **Organization:** Organize notes into notebooks or apply tags for easy searching and filtering.
*   **Attachments:** Attach images and other files to your notes.
*   **Search Functionality:** A powerful search to quickly find the information you need within your notes.

### Journaling
*   **Daily Entries:** Create daily journal entries to record your thoughts, experiences, and reflections.
*   **Mood Tracking:** Optionally track your mood for each journal entry.
*   **Photo Journaling:** Add photos to your journal entries to capture memories visually.
*   **Calendar View:** A calendar view to easily navigate through your past journal entries.
*   **Templates:** Use pre-defined templates for guided journaling, such as gratitude logs or daily reflections.

## 3. Architecture

Tempo is designed with a modern, scalable architecture to ensure a reliable and responsive user experience.

### 3.1. Backend

The backend is the core of the Tempo application, responsible for handling all business logic, data storage, and API requests.

*   **Language:** **Go (Golang)**
    Go is chosen for its high performance, concurrency support, and strong typing, making it an excellent choice for building scalable and efficient backend services. Go's ability to compile to a single binary simplifies deployment.

*   **Web Framework:** **Gin** or **Echo**
    To streamline development, a lightweight and high-performance web framework like Gin or Echo will be used. These frameworks provide robust routing, middleware support, and simplify the process of building a RESTful API.

*   **Database:** **PostgreSQL**
    PostgreSQL is a powerful, open-source object-relational database system known for its reliability, feature robustness, and performance. It provides a solid foundation for storing the application's data with a flexible schema.

*   **API:** **RESTful API**
    The backend will expose a RESTful API to allow the web and Android applications to communicate with the server. The API will use JSON for data exchange.

### 3.2. Frontend

#### 3.2.1. Web Application

The web application will provide a rich and interactive user experience for managing to-do lists, notes, and journal entries from any modern web browser.

*   **Framework:** **Nuxt.js**. The web application will be built using Nuxt.js, a powerful and intuitive framework based on Vue.js. Nuxt.js is chosen for its excellent developer experience, server-side rendering (SSR) capabilities for better performance and SEO, and its structured approach to building scalable single-page applications (SPAs).

#### 3.2.2. Android App

The Android app will offer a native experience for on-the-go productivity, ensuring the best possible performance and integration with the Android operating system.

*   **Language:** **Kotlin**. The application will be developed natively for Android using Kotlin, the modern, official language for Android development. This choice guarantees optimal performance, direct access to the latest platform APIs and features, and a superior user experience that feels seamless and responsive.

## 4. Database Schema

The PostgreSQL database will have a clear and normalized structure to efficiently store user data. Below is a simplified representation of the core tables:

```sql
-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- To-Do Lists Table
CREATE TABLE todo_lists (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- To-Do Items Table
CREATE TABLE todo_items (
    id SERIAL PRIMARY KEY,
    list_id INTEGER REFERENCES todo_lists(id) ON DELETE CASCADE,
    task TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    due_date DATE,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Notes Table
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Journal Entries Table
CREATE TABLE journal_entries (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    mood VARCHAR(50),
    entry_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## 5. API Endpoints

The RESTful API will provide the following endpoints for core functionalities (this is a non-exhaustive list):

### Users
*   `POST /api/users/register`: Register a new user.
*   `POST /api/users/login`: Authenticate a user and receive a token.

### To-Do Lists
*   `GET /api/lists`: Get all to-do lists for the authenticated user.
*   `POST /api/lists`: Create a new to-do list.
*   `GET /api/lists/{listId}`: Get a specific to-do list and its items.
*   `PUT /api/lists/{listId}`: Update a to-do list's title.
*   `DELETE /api/lists/{listId}`: Delete a to-do list.

### To-Do Items
*   `POST /api/lists/{listId}/items`: Create a new to-do item in a list.
*   `PUT /api/items/{itemId}`: Update a to-do item (e.g., mark as complete, change due date).
*   `DELETE /api/items/{itemId}`: Delete a to-do item.

### Notes
*   `GET /api/notes`: Get all notes for the authenticated user.
*   `POST /api/notes`: Create a new note.
*   `GET /api/notes/{noteId}`: Get a specific note.
*   `PUT /api/notes/{noteId}`: Update a note.
*   `DELETE /api/notes/{noteId}`: Delete a note.

### Journal Entries
*   `GET /api/journal`: Get all journal entries for the authenticated user.
*   `POST /api/journal`: Create a new journal entry.
*   `GET /api/journal/{entryId}`: Get a specific journal entry.
*   `PUT /api/journal/{entryId}`: Update a journal entry.
*   `DELETE /api/journal/{entryId}`: Delete a journal entry.

## 6. Deployment

*   **Backend:** The Go backend will be containerized using **Docker** and deployed on **Google Cloud Run**. This serverless platform will automatically scale the application based on traffic, providing a highly scalable and cost-effective solution.
*   **Web App:** The Nuxt.js frontend will be deployed on **Vercel**. Vercel is an ideal platform for Nuxt.js applications, offering seamless Git integration, automatic builds, and a global CDN for optimal performance.
*   **Android App:** The Android application will be packaged and distributed through the **Google Play Store**.