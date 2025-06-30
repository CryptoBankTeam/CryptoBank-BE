# Гибридная децентрализованная платформа для безопасной аренды криптоактивов с залогом через смарт-контракты в сети Ethereum (Sepolia).

##  Серверная логика 
-  **Отображение сделок и информации о пользователе в интерфейсе**
-  **Слушатель событий** смарт-контракта (LoanCreated, LoanRepaid и др.).
-  **Планировщик** для проверки просроченных сделок и синхронизаии данных при сбоях работы слушателя.
-  **Синхронизация** данных между Ethereum и PostgreSQL.
-  **JWT-аутентификация** и привязка кошельков к аккаунтам.
-  **Рейтинговая система** на основе истории транзакций.

##  Технологический стек
| Категория       | Технологии                                                                 |
|-----------------|----------------------------------------------------------------------------|
| Frontend        | React, Next.js, TypeScript, TailwindCSS                                    |
| Backend         | Go, Gin-gonic, Node.js, ether.js, JWT                                      |
| Web3            | ethers.js, MetaMask, Sepolia/ etherscan                                    |
| Тестирование    | Hardhat                                                                    |
| Деплой          | Hardhat                                                                    |

## Установка и запуск
1. Убедитесь, что установлены:
  - *Go 1.20+*
  - *PostgreSQL 16+*
2. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/CryptoBankTeam/CryptoBank-BE.git
   cd CryptoBank-FE
3. Настройте .env:
   ```bash
   DB_URL=postgres://user:pass@localhost:5432/cryptobank
   ALCHEMY_URL=https://eth-sepolia.g.alchemy.com/...
4. Запустите сервер
   ```bash
   go run main.go
5. Запустите слушатель событий
   ```bash
   cd event-sync
   node.js sync-loans.js
5. Запустите планировщик
   ```bash
   cd scripts
   node.js cronOverdue.js
