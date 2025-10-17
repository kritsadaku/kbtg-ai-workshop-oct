mod domain;
mod infrastructure;
mod application;
mod presentation;

use std::sync::Arc;
use utoipa::OpenApi;
use utoipa_swagger_ui::SwaggerUi;
use sqlx::SqlitePool;

use domain::{User, CreateUserRequest, UpdateUserRequest, Transfer, CreateTransferRequest, TransferCreateResponse, TransferGetResponse, TransferListResponse};
use infrastructure::{SqliteUserRepository, SqliteTransferRepository, SqlitePointLedgerRepository};
use application::{UserService, TransferService};
use presentation::{create_routes, AppState, ErrorResponse, ListUsersResponse};

#[derive(OpenApi)]
#[openapi(
    paths(
        presentation::handlers::hello_world,
        presentation::handlers::get_user,
        presentation::handlers::list_users,
        presentation::handlers::create_user,
        presentation::handlers::update_user,
        presentation::handlers::delete_user,
        presentation::transfer_handlers::create_transfer,
        presentation::transfer_handlers::get_transfer,
        presentation::transfer_handlers::list_transfers,
    ),
    components(
        schemas(User, CreateUserRequest, UpdateUserRequest, Transfer, CreateTransferRequest, TransferCreateResponse, TransferGetResponse, TransferListResponse, ErrorResponse, ListUsersResponse)
    ),
    tags(
        (name = "simple-app", description = "Clean Architecture API with User Management and SQLite")
    )
)]
struct ApiDoc;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Database setup
    let database_url = "sqlite:users.db";
    
    // Create database file if it doesn't exist
    if !std::path::Path::new("users.db").exists() {
        println!("📂 Creating SQLite database file: users.db");
        std::fs::File::create("users.db")?;
    }
    
    let pool = SqlitePool::connect(database_url).await?;
    
    // Infrastructure layer - Repositories
    let user_repository = Arc::new(SqliteUserRepository::new(pool.clone()));
    let transfer_repository = Arc::new(SqliteTransferRepository::new(pool.clone()));
    let point_ledger_repository = Arc::new(SqlitePointLedgerRepository::new(pool.clone()));
    
    // Initialize database tables
    user_repository.init_database().await?;
    transfer_repository.init_database().await?;
    point_ledger_repository.init_database().await?;
    
    // Application layer - Services
    let user_service = UserService::new(user_repository.clone());
    let transfer_service = TransferService::new(
        transfer_repository,
        point_ledger_repository,
        user_repository,
    );
    
    // Application state
    let app_state = AppState { 
        user_service,
        transfer_service,
    };
    
    // Presentation layer - Routes
    let app = create_routes()
        .merge(SwaggerUi::new("/swagger-ui").url("/api-docs/openapi.json", ApiDoc::openapi()))
        .with_state(app_state);

    println!("🚀 Server running on http://0.0.0.0:3000");
    println!("📚 Swagger UI available at http://0.0.0.0:3000/swagger-ui");
    println!("💾 SQLite database: users.db");
    println!("🔗 API Endpoints:");
    println!("   GET    /");
    println!("   GET    /users?limit=10&offset=0");
    println!("   POST   /users");
    println!("   GET    /users/{{id}}");
    println!("   PUT    /users/{{id}}");
    println!("   DELETE /users/{{id}}");
    println!("   POST   /transfers");
    println!("   GET    /transfers?userId={{userId}}&page=1&pageSize=20");
    println!("   GET    /transfers/{{id}}");
    println!("");
    println!("📊 Transfer API Features:");
    println!("   - Point transfer between users");
    println!("   - Idempotency key for duplicate protection");
    println!("   - Point ledger for audit trail");
    println!("   - Automatic balance management");

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await?;
    axum::serve(listener, app).await?;
    
    Ok(())
}
