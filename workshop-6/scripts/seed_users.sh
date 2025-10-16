#!/bin/bash

# Script สำหรับ seed ข้อมูล users จาก CSV file เข้า SQLite database
# Usage: ./seed_users.sh [csv_file] [database_file]

set -e  # Exit on any error

# Default values
CSV_FILE="${1:-./users_data.csv}"
DB_FILE="${2:-../users.db}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🌱 Starting user data seeding process...${NC}"

# Check if CSV file exists
if [ ! -f "$CSV_FILE" ]; then
    echo -e "${RED}❌ Error: CSV file '$CSV_FILE' not found!${NC}"
    echo -e "${YELLOW}💡 Usage: $0 [csv_file] [database_file]${NC}"
    echo -e "${YELLOW}💡 Example: $0 ./users_data.csv ../users.db${NC}"
    exit 1
fi

# Check if database file exists
if [ ! -f "$DB_FILE" ]; then
    echo -e "${RED}❌ Error: Database file '$DB_FILE' not found!${NC}"
    echo -e "${YELLOW}💡 Please make sure the database is initialized first by running the Go application${NC}"
    exit 1
fi

# Check if sqlite3 is installed
if ! command -v sqlite3 &> /dev/null; then
    echo -e "${RED}❌ Error: sqlite3 is not installed!${NC}"
    echo -e "${YELLOW}💡 Please install sqlite3 first:${NC}"
    echo -e "${YELLOW}   macOS: brew install sqlite${NC}"
    echo -e "${YELLOW}   Ubuntu: sudo apt-get install sqlite3${NC}"
    exit 1
fi

echo -e "${BLUE}📄 CSV file: $CSV_FILE${NC}"
echo -e "${BLUE}🗄️  Database: $DB_FILE${NC}"

# Count existing users before seeding
EXISTING_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM users;")
echo -e "${BLUE}📊 Existing users in database: $EXISTING_COUNT${NC}"

# Create temporary SQL file
TEMP_SQL=$(mktemp)
echo -e "${BLUE}📝 Creating temporary SQL file: $TEMP_SQL${NC}"

# Generate SQL INSERT statements from CSV
echo "-- Auto-generated SQL for seeding users" > "$TEMP_SQL"
echo "BEGIN TRANSACTION;" >> "$TEMP_SQL"

# Skip header line and process each data line
tail -n +2 "$CSV_FILE" | while IFS=',' read -r first_name last_name phone email member_since membership_level points; do
    # Get current timestamp
    current_time=$(date '+%Y-%m-%d %H:%M:%S')
    
    # Escape single quotes in data
    first_name=$(echo "$first_name" | sed "s/'/''/g")
    last_name=$(echo "$last_name" | sed "s/'/''/g")
    phone=$(echo "$phone" | sed "s/'/''/g")
    email=$(echo "$email" | sed "s/'/''/g")
    member_since=$(echo "$member_since" | sed "s/'/''/g")
    membership_level=$(echo "$membership_level" | sed "s/'/''/g")
    
    # Generate INSERT statement
    cat >> "$TEMP_SQL" << EOF
INSERT OR IGNORE INTO users (first_name, last_name, phone, email, member_since, membership_level, points, created_at, updated_at)
VALUES ('$first_name', '$last_name', '$phone', '$email', '$member_since', '$membership_level', $points, '$current_time', '$current_time');
EOF
done

echo "COMMIT;" >> "$TEMP_SQL"

echo -e "${YELLOW}🔄 Processing CSV data...${NC}"

# Count lines to process (excluding header)
TOTAL_LINES=$(($(wc -l < "$CSV_FILE") - 1))
echo -e "${BLUE}📈 Found $TOTAL_LINES users to process${NC}"

# Execute SQL file
echo -e "${YELLOW}💾 Inserting data into database...${NC}"
if sqlite3 "$DB_FILE" < "$TEMP_SQL"; then
    echo -e "${GREEN}✅ Data insertion completed successfully!${NC}"
else
    echo -e "${RED}❌ Error occurred during data insertion!${NC}"
    rm -f "$TEMP_SQL"
    exit 1
fi

# Clean up temporary file
rm -f "$TEMP_SQL"

# Count users after seeding
NEW_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM users;")
ADDED_COUNT=$((NEW_COUNT - EXISTING_COUNT))

echo -e "${GREEN}🎉 Seeding completed successfully!${NC}"
echo -e "${BLUE}📊 Statistics:${NC}"
echo -e "${BLUE}   • Users before seeding: $EXISTING_COUNT${NC}"
echo -e "${BLUE}   • Users after seeding: $NEW_COUNT${NC}"
echo -e "${GREEN}   • New users added: $ADDED_COUNT${NC}"

# Show sample of inserted data
echo -e "${BLUE}📋 Sample of recent users:${NC}"
sqlite3 -header -column "$DB_FILE" "SELECT id, first_name, last_name, email, membership_level, points FROM users ORDER BY id DESC LIMIT 5;"

echo -e "${GREEN}✨ Seeding process completed!${NC}"