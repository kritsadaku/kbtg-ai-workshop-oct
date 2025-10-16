#!/bin/bash

# Script สำหรับนับจำนวน users ใน SQLite database และแสดงสถิติต่างๆ
# Usage: ./count_users.sh [database_file]

set -e  # Exit on any error

# Default values
DB_FILE="${1:-../users.db}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}📊 User Database Statistics${NC}"
echo -e "${BLUE}=========================${NC}"

# Check if database file exists
if [ ! -f "$DB_FILE" ]; then
    echo -e "${RED}❌ Error: Database file '$DB_FILE' not found!${NC}"
    echo -e "${YELLOW}💡 Usage: $0 [database_file]${NC}"
    echo -e "${YELLOW}💡 Example: $0 ../users.db${NC}"
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

echo -e "${BLUE}🗄️  Database: $DB_FILE${NC}"
echo

# Basic count
echo -e "${CYAN}📈 Total User Count${NC}"
echo -e "${CYAN}==================${NC}"
TOTAL_USERS=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM users;")
echo -e "${GREEN}👥 Total Users: $TOTAL_USERS${NC}"
echo

# Membership level breakdown
echo -e "${CYAN}🏆 Membership Level Breakdown${NC}"
echo -e "${CYAN}=============================${NC}"
sqlite3 -header -column "$DB_FILE" "
SELECT 
    membership_level as 'Membership Level',
    COUNT(*) as 'Count',
    ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM users), 2) || '%' as 'Percentage'
FROM users 
GROUP BY membership_level 
ORDER BY 
    CASE membership_level 
        WHEN 'Bronze' THEN 1 
        WHEN 'Silver' THEN 2 
        WHEN 'Gold' THEN 3 
        WHEN 'Platinum' THEN 4 
        ELSE 5 
    END;
"
echo

# Points statistics
echo -e "${CYAN}💰 Points Statistics${NC}"
echo -e "${CYAN}===================${NC}"
sqlite3 -header -column "$DB_FILE" "
SELECT 
    'Total Points' as 'Metric',
    SUM(points) as 'Value'
FROM users
UNION ALL
SELECT 
    'Average Points' as 'Metric',
    ROUND(AVG(points), 2) as 'Value'
FROM users
UNION ALL
SELECT 
    'Minimum Points' as 'Metric',
    MIN(points) as 'Value'
FROM users
UNION ALL
SELECT 
    'Maximum Points' as 'Metric',
    MAX(points) as 'Value'
FROM users;
"
echo

# Registration timeline
echo -e "${CYAN}📅 Registration Timeline${NC}"
echo -e "${CYAN}========================${NC}"
sqlite3 -header -column "$DB_FILE" "
SELECT 
    strftime('%Y-%m', member_since) as 'Month',
    COUNT(*) as 'New Users'
FROM users 
GROUP BY strftime('%Y-%m', member_since)
ORDER BY strftime('%Y-%m', member_since);
"
echo

# Top 10 users by points
echo -e "${CYAN}🌟 Top 10 Users by Points${NC}"
echo -e "${CYAN}=========================${NC}"
sqlite3 -header -column "$DB_FILE" "
SELECT 
    id as 'ID',
    first_name || ' ' || last_name as 'Name',
    email as 'Email',
    membership_level as 'Level',
    points as 'Points'
FROM users 
ORDER BY points DESC 
LIMIT 10;
"
echo

# Recent registrations
echo -e "${CYAN}🆕 Recent Registrations (Last 5)${NC}"
echo -e "${CYAN}================================${NC}"
sqlite3 -header -column "$DB_FILE" "
SELECT 
    id as 'ID',
    first_name || ' ' || last_name as 'Name',
    email as 'Email',
    membership_level as 'Level',
    points as 'Points',
    created_at as 'Registered'
FROM users 
ORDER BY created_at DESC 
LIMIT 5;
"
echo

# Database size information
echo -e "${CYAN}💾 Database Information${NC}"
echo -e "${CYAN}======================${NC}"
DB_SIZE=$(ls -lh "$DB_FILE" | awk '{print $5}')
echo -e "${BLUE}📁 Database file size: $DB_SIZE${NC}"

# Table record counts
echo -e "${BLUE}📋 Table Record Counts:${NC}"
USERS_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM users;")
TRANSFERS_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM transfers;" 2>/dev/null || echo "0")
LEDGER_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM point_ledger;" 2>/dev/null || echo "0")

echo -e "${BLUE}   • Users: $USERS_COUNT${NC}"
echo -e "${BLUE}   • Transfers: $TRANSFERS_COUNT${NC}"
echo -e "${BLUE}   • Point Ledger: $LEDGER_COUNT${NC}"
echo

# Summary
echo -e "${GREEN}✨ Summary${NC}"
echo -e "${GREEN}==========${NC}"
if [ "$TOTAL_USERS" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  No users found in the database${NC}"
    echo -e "${YELLOW}💡 Run the seed script to add sample data: ./seed_users.sh${NC}"
elif [ "$TOTAL_USERS" -lt 10 ]; then
    echo -e "${YELLOW}📊 Database has a small number of users ($TOTAL_USERS)${NC}"
    echo -e "${YELLOW}💡 Consider adding more sample data for testing${NC}"
else
    echo -e "${GREEN}✅ Database is well populated with $TOTAL_USERS users${NC}"
fi

echo -e "${GREEN}🎉 Statistics generation completed!${NC}"