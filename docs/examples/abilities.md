# Примеры реализации Abilities

Здесь собраны реальные примеры Abilities для различных сценариев тестирования.

## 🗄️ Database Ability (PostgreSQL)

Полная реализация для работы с PostgreSQL базой данных.

### Интерфейс

```go
package database

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    
    "github.com/nchursin/serenity-go/serenity/abilities"
)

// DatabaseAbility - способность для работы с базами данных
type DatabaseAbility interface {
    abilities.Ability
    
    // Connection management
    Connect(dsn string) error
    Disconnect() error
    Ping() error
    
    // Query operations
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
    Execute(query string, args ...interface{}) (sql.Result, error)
    
    // Transaction management
    BeginTx() (*sql.Tx, error)
    
    // State information
    LastQuery() string
    LastError() error
    IsConnected() bool
}
```

### Реализация

```go
type databaseAbility struct {
    db         *sql.DB
    lastQuery  string
    lastError  error
    dsn        string
    mutex      sync.RWMutex
}

// ConnectToPostgreSQL - создает способность для подключения к PostgreSQL
func ConnectToPostgreSQL(dsn string) DatabaseAbility {
    return &databaseAbility{
        dsn: dsn,
    }
}

func (d *databaseAbility) Connect(dsn string) error {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    
    if dsn != "" {
        d.dsn = dsn
    }
    
    db, err := sql.Open("postgres", d.dsn)
    if err != nil {
        d.lastError = fmt.Errorf("failed to open database: %w", err)
        return d.lastError
    }
    
    // Проверяем соединение
    if err := db.Ping(); err != nil {
        d.lastError = fmt.Errorf("failed to ping database: %w", err)
        db.Close()
        return d.lastError
    }
    
    d.db = db
    d.lastError = nil
    return nil
}

func (d *databaseAbility) Query(query string, args ...interface{}) (*sql.Rows, error) {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    
    if d.db == nil {
        err := fmt.Errorf("database not connected")
        d.lastError = err
        d.lastQuery = query
        return nil, err
    }
    
    d.lastQuery = query
    rows, err := d.db.Query(query, args...)
    d.lastError = err
    return rows, err
}

func (d *databaseAbility) Execute(query string, args ...interface{}) (sql.Result, error) {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    
    if d.db == nil {
        err := fmt.Errorf("database not connected")
        d.lastError = err
        d.lastQuery = query
        return nil, err
    }
    
    d.lastQuery = query
    result, err := d.db.Exec(query, args...)
    d.lastError = err
    return result, err
}

func (d *databaseAbility) Disconnect() error {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    
    if d.db != nil {
        err := d.db.Close()
        d.db = nil
        d.lastError = err
        return err
    }
    return nil
}
```

### Activities

```go
// CreateTableActivity - создание таблицы
type CreateTableActivity struct {
    tableName string
    schema    string
}

func CreateTable(tableName, schema string) *CreateTableActivity {
    return &CreateTableActivity{
        tableName: tableName,
        schema:    schema,
    }
}

func (c *CreateTableActivity) PerformAs(actor core.Actor) error {
    ability, err := actor.AbilityTo(&databaseAbility{})
    if err != nil {
        return fmt.Errorf("actor does not have database ability: %w", err)
    }
    
    db := ability.(DatabaseAbility)
    _, err = db.Execute(c.schema)
    return err
}

func (c *CreateTableActivity) Description() string {
    return fmt.Sprintf("creates table: %s", c.tableName)
}

// InsertDataActivity - вставка данных
type InsertDataActivity struct {
    table string
    data  map[string]interface{}
}

func InsertInto(table string, data map[string]interface{}) *InsertDataActivity {
    return &InsertDataActivity{table: table, data: data}
}

func (i *InsertDataActivity) PerformAs(actor core.Actor) error {
    ability, err := actor.AbilityTo(&databaseAbility{})
    if err != nil {
        return fmt.Errorf("actor does not have database ability: %w", err)
    }
    
    db := ability.(DatabaseAbility)
    
    // Build INSERT query
    columns := make([]string, 0, len(i.data))
    placeholders := make([]string, 0, len(i.data))
    values := make([]interface{}, 0, len(i.data))
    
    for column, value := range i.data {
        columns = append(columns, column)
        placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
        values = append(values, value)
    }
    
    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
        i.table,
        strings.Join(columns, ", "),
        strings.Join(placeholders, ", "),
    )
    
    _, err = db.Execute(query, values...)
    return err
}
```

### Использование

```go
func TestDatabaseOperations(t *testing.T) {
    actor := core.NewActor("DBAdmin").WhoCan(
        database.ConnectToPostgreSQL("postgres://user:pass@localhost/testdb?sslmode=disable"),
    )
    
    err := actor.AttemptsTo(
        // Подключаемся к базе
        core.Do("connects to database", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&databaseAbility{})
            return ability.(DatabaseAbility).Connect("")
        }),
        
        // Создаем таблицу
        database.CreateTable("users", `
            CREATE TABLE users (
                id SERIAL PRIMARY KEY,
                name VARCHAR(100) NOT NULL,
                email VARCHAR(255) UNIQUE NOT NULL,
                created_at TIMESTAMP DEFAULT NOW()
            )
        `),
        
        // Вставляем данные
        database.InsertInto("users", map[string]interface{}{
            "name":  "John Doe",
            "email": "john@example.com",
        }),
        
        // Проверяем результат
        ensure.That(database.RowCount("users"), expectations.Equals(1)),
    )
    
    require.NoError(t, err)
    
    // Очистка
    cleanupActor := core.NewActor("DBCleaner").WhoCan(
        database.ConnectToPostgreSQL("postgres://user:pass@localhost/testdb?sslmode=disable"),
    )
    
    cleanupErr := cleanupActor.AttemptsTo(
        core.Do("connects to database", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&databaseAbility{})
            return ability.(DatabaseAbility).Connect("")
        }),
        database.DropTable("users"),
    )
    
    require.NoError(t, cleanupErr)
}
```

## 🗂️ FileSystem Ability

Расширенная версия файловой системы с backup и версионированием.

### Интерфейс

```go
package filesystem

import (
    "io/fs"
    "time"
)

// FileSystemAbility - расширенная способность для работы с файлами
type FileSystemAbility interface {
    abilities.Ability
    
    // Basic operations
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, data []byte, perm fs.FileMode) error
    DeleteFile(path string) error
    Exists(path string) bool
    
    // Directory operations
    CreateDir(path string, perm fs.FileMode) error
    ListDir(path string) ([]fs.DirEntry, error)
    
    // Advanced features
    BackupFile(path string) (string, error)
    RestoreFile(backupPath string) error
    GetFileSize(path string) int64
    GetFileModTime(path string) time.Time
    
    // Working directory
    SetWorkingDirectory(dir string) error
    GetWorkingDirectory() string
    
    // State
    LastOperation() string
}
```

### Реализация

```go
type fileSystemAbility struct {
    workingDir   string
    backupDir    string
    lastOp       string
    backups      map[string]string // original -> backup path
    mutex        sync.RWMutex
}

func ManageFileSystem() FileSystemAbility {
    return &fileSystemAbility{
        workingDir: ".",
        backupDir:  ".backups",
        backups:    make(map[string]string),
    }
}

func ManageFileSystemIn(directory string) FileSystemAbility {
    abs, _ := filepath.Abs(directory)
    return &fileSystemAbility{
        workingDir: abs,
        backupDir:  filepath.Join(abs, ".backups"),
        backups:    make(map[string]string),
    }
}

func (f *fileSystemAbility) WriteFile(path string, data []byte, perm fs.FileMode) error {
    f.mutex.Lock()
    defer f.mutex.Unlock()
    
    fullPath := filepath.Join(f.workingDir, path)
    
    // Создаем директорию если нужно
    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        f.lastOp = fmt.Sprintf("mkdir error: %s", path)
        return fmt.Errorf("failed to create directory: %w", err)
    }
    
    if err := os.WriteFile(fullPath, data, perm); err != nil {
        f.lastOp = fmt.Sprintf("write error: %s", path)
        return fmt.Errorf("failed to write file: %w", err)
    }
    
    f.lastOp = fmt.Sprintf("write: %s", path)
    return nil
}

func (f *fileSystemAbility) BackupFile(path string) (string, error) {
    f.mutex.Lock()
    defer f.mutex.Unlock()
    
    fullPath := filepath.Join(f.workingDir, path)
    
    // Создаем backup директорию
    if err := os.MkdirAll(f.backupDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create backup directory: %w", err)
    }
    
    // Генерируем имя backup файла
    timestamp := time.Now().Format("20060102-150405")
    backupName := fmt.Sprintf("%s_%s_%s", 
        filepath.Base(path), 
        strings.TrimSuffix(filepath.Ext(path), "."), 
        timestamp,
    )
    backupPath := filepath.Join(f.backupDir, backupName)
    
    // Копируем файл
    if err := copyFile(fullPath, backupPath); err != nil {
        return "", fmt.Errorf("failed to backup file: %w", err)
    }
    
    f.backups[path] = backupPath
    f.lastOp = fmt.Sprintf("backup: %s -> %s", path, backupPath)
    return backupPath, nil
}

func (f *fileSystemAbility) RestoreFile(backupPath string) error {
    f.mutex.Lock()
    defer f.mutex.Unlock()
    
    // Находим оригинальный путь
    var originalPath string
    for orig, backup := range f.backups {
        if backup == backupPath {
            originalPath = orig
            break
        }
    }
    
    if originalPath == "" {
        return fmt.Errorf("backup not found: %s", backupPath)
    }
    
    fullPath := filepath.Join(f.workingDir, originalPath)
    
    if err := copyFile(backupPath, fullPath); err != nil {
        f.lastOp = fmt.Sprintf("restore error: %s", originalPath)
        return fmt.Errorf("failed to restore file: %w", err)
    }
    
    f.lastOp = fmt.Sprintf("restore: %s", originalPath)
    return nil
}

func copyFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()
    
    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()
    
    _, err = io.Copy(destination, source)
    return err
}
```

### Использование

```go
func TestFileSystemWithBackup(t *testing.T) {
    tempDir := t.TempDir()
    
    actor := core.NewActor("FileEditor").WhoCan(
        filesystem.ManageFileSystemIn(tempDir),
    )
    
    originalContent := "original content"
    modifiedContent := "modified content"
    
    err := actor.AttemptsTo(
        // Создаем оригинальный файл
        core.Do("creates original file", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&fileSystemAbility{})
            return ability.(FileSystemAbility).WriteFile(
                "important.txt", 
                []byte(originalContent), 
                0644,
            )
        }),
        
        // Делаем backup
        filesystem.BackupUpFile("important.txt"),
        
        // Модифицируем файл
        core.Do("modifies file", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&fileSystemAbility{})
            return ability.(FileSystemAbility).WriteFile(
                "important.txt", 
                []byte(modifiedContent), 
                0644,
            )
        }),
        
        // Проверяем модификацию
        ensure.That(filesystem.FileContent("important.txt"), expectations.Equals(modifiedContent)),
        
        // Восстанавливаем из backup
        filesystem.RestoreLastBackup("important.txt"),
        
        // Проверяем восстановление
        ensure.That(filesystem.FileContent("important.txt"), expectations.Equals(originalContent)),
    )
    
    require.NoError(t, err)
}
```

## 🔌 WebSocket Ability

Способность для работы с WebSocket соединениями.

### Интерфейс

```go
package websocket

import (
    "time"
    "github.com/gorilla/websocket"
)

// WebSocketAbility - способность для работы с WebSocket
type WebSocketAbility interface {
    abilities.Ability
    
    // Connection management
    Connect(url string, header http.Header) error
    Disconnect() error
    IsConnected() bool
    
    // Messaging
    Send(message []byte) error
    SendJSON(v interface{}) error
    Receive(timeout time.Duration) ([]byte, error)
    ReceiveJSON(v interface{}, timeout time.Duration) error
    
    // State
    LastMessage() []byte
    ConnectionDuration() time.Duration
    MessageCount() int
}
```

### Реализация

```go
type webSocketAbility struct {
    conn           *websocket.Conn
    lastMessage    []byte
    connectTime    time.Time
    messageCount   int
    mutex          sync.RWMutex
    dialer         websocket.Dialer
    pingInterval   time.Duration
    pongWait       time.Duration
}

func ConnectToWebSocket() WebSocketAbility {
    return &webSocketAbility{
        dialer:       *websocket.DefaultDialer,
        pingInterval: 30 * time.Second,
        pongWait:     60 * time.Second,
    }
}

func ConnectToWebSocketWithTimeout(timeout time.Duration) WebSocketAbility {
    dialer := websocket.DefaultDialer
    dialer.HandshakeTimeout = timeout
    
    return &webSocketAbility{
        dialer:       *dialer,
        pingInterval: 30 * time.Second,
        pongWait:     60 * time.Second,
    }
}

func (w *webSocketAbility) Connect(url string, header http.Header) error {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    conn, resp, err := w.dialer.Dial(url, header)
    if err != nil {
        return fmt.Errorf("failed to connect to websocket: %w", err)
    }
    
    if resp != nil && resp.Body != nil {
        resp.Body.Close()
    }
    
    w.conn = conn
    w.connectTime = time.Now()
    w.messageCount = 0
    
    // Запускаем ping/pong handler
    go w.startPingPong()
    
    return nil
}

func (w *webSocketAbility) Send(message []byte) error {
    w.mutex.RLock()
    defer w.mutex.RUnlock()
    
    if w.conn == nil {
        return fmt.Errorf("not connected to websocket")
    }
    
    err := w.conn.WriteMessage(websocket.TextMessage, message)
    if err == nil {
        w.messageCount++
    }
    return err
}

func (w *webSocketAbility) Receive(timeout time.Duration) ([]byte, error) {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    if w.conn == nil {
        return nil, fmt.Errorf("not connected to websocket")
    }
    
    // Set read deadline
    if timeout > 0 {
        err := w.conn.SetReadDeadline(time.Now().Add(timeout))
        if err != nil {
            return nil, fmt.Errorf("failed to set read deadline: %w", err)
        }
    }
    
    messageType, message, err := w.conn.ReadMessage()
    if err != nil {
        return nil, fmt.Errorf("failed to read message: %w", err)
    }
    
    if messageType == websocket.TextMessage {
        w.lastMessage = message
        w.messageCount++
    }
    
    return message, nil
}

func (w *webSocketAbility) startPingPong() {
    ticker := time.NewTicker(w.pingInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            w.mutex.RLock()
            conn := w.conn
            w.mutex.RUnlock()
            
            if conn == nil {
                return
            }
            
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                conn.Close()
                return
            }
        }
    }
}
```

### Использование

```go
func TestWebSocketChat(t *testing.T) {
    // Запускаем тестовый WebSocket сервер
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        upgrader := websocket.Upgrader{}
        conn, _ := upgrader.Upgrade(w, r, nil)
        
        defer conn.Close()
        
        for {
            messageType, message, err := conn.ReadMessage()
            if err != nil {
                break
            }
            
            // Echo message back
            conn.WriteMessage(messageType, message)
        }
    }))
    defer server.Close()
    
    // Конвертируем HTTP URL в WebSocket URL
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
    
    actor := core.NewActor("WebSocketClient").WhoCan(
        websocket.ConnectToWebSocket(),
    )
    
    err := actor.AttemptsTo(
        // Подключаемся к WebSocket
        core.Do("connects to websocket", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&webSocketAbility{})
            return ability.(WebSocketAbility).Connect(wsURL, nil)
        }),
        
        // Отправляем сообщение
        core.Do("sends message", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&webSocketAbility{})
            return ability.(WebSocketAbility).Send([]byte("Hello WebSocket!"))
        }),
        
        // Получаем ответ
        core.Do("receives response", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&webSocketAbility{})
            _, err := ability.(WebSocketAbility).Receive(5 * time.Second)
            return err
        }),
        
        // Проверяем полученное сообщение
        ensure.That(websocket.LastMessage(), expectations.Equals([]byte("Hello WebSocket!"))),
        ensure.That(websocket.MessageCount(), expectations.Equals(2)), // отправлено + получено
    )
    
    require.NoError(t, err)
}
```

## 📊 Redis Ability

Способность для работы с Redis.

### Интерфейс

```go
package redis

import "github.com/go-redis/redis/v8"

// RedisAbility - способность для работы с Redis
type RedisAbility interface {
    abilities.Ability
    
    // Connection
    Connect(addr string, options *redis.Options) error
    Disconnect() error
    Ping() error
    
    // Basic operations
    Set(key string, value interface{}, expiration time.Duration) error
    Get(key string) (string, error)
    Del(keys ...string) error
    Exists(keys ...string) (int64, error)
    
    // Hash operations
    HSet(key string, values ...interface{}) error
    HGet(key, field string) (string, error)
    HGetAll(key string) (map[string]string, error)
    
    // List operations
    LPush(key string, values ...interface{}) error
    RPop(key string) (string, error)
    LRange(key string, start, stop int64) ([]string, error)
    
    // State
    LastCommand() string
    ConnectionInfo() string
}
```

### Пример использования

```go
func TestRedisOperations(t *testing.T) {
    // Предполагаем, что у вас запущен Redis на localhost:6379
    actor := core.NewActor("RedisUser").WhoCan(
        redis.ConnectToRedis("localhost:6379"),
    )
    
    err := actor.AttemptsTo(
        // Подключаемся к Redis
        core.Do("connects to redis", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&redisAbility{})
            return ability.(RedisAbility).Connect("", &redis.Options{
                Addr: "localhost:6379",
            })
        }),
        
        // Устанавливаем значение
        core.Do("sets key-value", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&redisAbility{})
            return ability.(RedisAbility).Set("test:key", "test-value", 0)
        }),
        
        // Получаем значение
        core.Do("gets value", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&redisAbility{})
            val, err := ability.(RedisAbility).Get("test:key")
            if err != nil {
                return err
            }
            // Сохраняем значение для проверки
            return nil
        }),
        
        // Проверяем существование ключа
        ensure.That(redis.KeyExists("test:key"), expectations.IsTrue()),
        ensure.That(redis.StringValue("test:key"), expectations.Equals("test-value")),
        
        // Удаляем ключ
        core.Do("deletes key", func(actor core.Actor) error {
            ability, _ := actor.AbilityTo(&redisAbility{})
            return ability.(RedisAbility).Del("test:key")
        }),
        
        // Проверяем удаление
        ensure.That(redis.KeyExists("test:key"), expectations.IsFalse()),
    )
    
    require.NoError(t, err)
}
```

---

Эти примеры показывают различные подходы к созданию Abilities:

1. **Database Ability** - классический пример с connection management
2. **FileSystem Ability** - продвинутая версия с backup функциями
3. **WebSocket Ability** - работа с real-time коммуникациями
4. **Redis Ability** - интеграция с популярным in-memory storage

Каждый пример следует паттернам, описанным в основной [инструкции по созданию Abilities](../abilities.md), и может быть адаптирован под ваши конкретные нужды.