# 🎟️ Tiketo API
Api for ticket order system, integrated with MIdtrans as payment gateway.

# 🚀 Features
- User authentication and authorization
- Ticket management, like create, update, delete, search.
- Order ticket
- Caching with redis

# 🛠️ Tech Stack
- Go
- PostgreSQL
- Gorm
- Redis
- Echo

# ⚙️ How to Use
### 🐳 Docker
1. Clone this project
```
git clone https://github.com/insanXYZ/tiketo
```
2. Set all environment on ``docker-compose.yaml``
3. Build/create compose 
```
docker compose create
```
4. Run app
```
docker compose start
```
### 💻 Local
1. Clone this project
```
git clone https://github.com/insanXYZ/tiketo
```
2. Rename ``.env.example`` to ``.env`` and set all configuration
3. Run app
```
make run
```