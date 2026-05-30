# Database Schema Relations

This document captures the entity-relationship structure of the Anugerah Jaya Farm Volare database, derived from the GORM models in `internal/entity/`. Relationships are drawn from `foreignKey` / `references` tags. Cardinality on the FK side reflects whether the FK column is `NOT NULL` (mandatory) or nullable (optional).

```mermaid
erDiagram
    %% ============================================================
    %% Core: tenancy + identity
    %% ============================================================
    Location {
        uint64 Id PK
        string Name
        float64 Latitude
        float64 Longitude
    }

    Role {
        uint64 Id PK
        string Name
    }

    User {
        uuid Id PK
        string Username
        string Email
        uint64 LocationId FK
        uint64 RoleId FK
        string Name
        string PhoneNumber
        int SalaryInterval
        decimal Salary
    }

    Location ||--o{ User : "has users"
    Role    ||--o{ User : "assigns role to"

    %% ============================================================
    %% Placements (user assigned to a cage/store/warehouse)
    %% ============================================================
    CagePlacement {
        uuid UserId PK,FK
        uint64 CageId PK,FK
    }
    StorePlacement {
        uuid UserId PK,FK
        uint64 StoreId FK
    }
    WarehousePlacement {
        uuid UserId PK,FK
        uint64 WarehouseId FK
    }

    User      ||--o{ CagePlacement      : "assigned"
    Cage      ||--o{ CagePlacement      : "staffed by"
    User      ||--o{ StorePlacement     : "assigned"
    Store     ||--o{ StorePlacement     : "staffed by"
    User      ||--o{ WarehousePlacement : "assigned"
    Warehouse ||--o{ WarehousePlacement : "staffed by"

    %% ============================================================
    %% Cage + chicken population
    %% ============================================================
    Cage {
        uint64 Id PK
        uint64 LocationId FK
        string Name
        uint64 Capacity
        int ChickenCategory
        bool IsUsed
    }
    ChickenCage {
        uint64 Id PK
        uint64 CageId FK
        int64 ChickenProcurementId FK "nullable"
        uint64 TotalChicken
        int LatestChickenAgeVaccineRoutine
        bool IsNeedRoutineVaccine
        bool IsNeedFeed
    }
    ChickenMonitoring {
        uint64 Id PK
        uint64 ChickenCageId FK
        uint64 TotalChicken
        uint64 TotalDeathChicken
        float64 MortalityRate
        uint64 TotalSickChicken
        float64 TotalFeed
    }
    ChickenHealthItem {
        uint64 Id PK
        string Name
        int Type
        int64 ChickenAge
    }
    ChickenHealthMonitoring {
        uint64 Id PK
        uint64 ChickenCageId FK
        string HealthItemName
        int Type
        float64 Dose
        string Unit
        uint64 ChickenAge
    }
    ChickenPerformance {
        uint64 Id PK
        uint64 ChickenCageId FK
        uint64 LocationId FK
        string CageName
        int ChickenCategory
        uint64 ChickenAge
        uint64 TotalChicken
        uint64 TotalEgg
        float64 FCR
        float64 HDP
        float64 MortalityRate
        int Productivity
        bool IsGetFeed
    }

    Location    ||--o{ Cage                    : "hosts"
    Cage        ||--o{ ChickenCage             : "has population"
    ChickenCage ||--o{ ChickenMonitoring       : "monitored by"
    ChickenCage ||--o{ ChickenHealthMonitoring : "health monitored by"
    ChickenCage ||--o{ ChickenPerformance      : "scored by"

    %% ============================================================
    %% Chicken procurement
    %% ============================================================
    ChickenProcurementDraft {
        uint64 Id PK
        uint64 CageId FK
        uint64 SupplierId FK
        uint64 Quantity
        decimal TotalPrice
    }
    ChickenProcurement {
        uint64 Id PK
        uint64 CageId FK
        uint64 SupplierId FK
        uint64 Quantity
        int64 ReceiveQuantity
        decimal TotalPrice
        int PaymentStatus
        int Status
        int PaymentType
        bool IsArrived
        time EstimationArrivalDate
        time DeadlinePaymentDate
        time PaidDate
        uuid CreatedBy FK
    }
    ChickenProcurementPayment {
        uint64 Id PK
        uint64 ChickenProcurementId FK
        decimal Nominal
        time PaymentDate
        uuid CreatedBy FK
    }

    Cage     ||--o{ ChickenProcurementDraft   : "drafts"
    Supplier ||--o{ ChickenProcurementDraft   : "from"
    Cage     ||--o{ ChickenProcurement        : "procures"
    Supplier ||--o{ ChickenProcurement        : "supplies"
    User     ||--o{ ChickenProcurement        : "created"
    ChickenProcurement ||--o{ ChickenProcurementPayment : "paid by"
    ChickenProcurement ||--o{ ChickenCage                : "fills"
    User     ||--o{ ChickenProcurementPayment : "created"

    %% ============================================================
    %% Cage feed (recipe + stock at cage)
    %% ============================================================
    CageFeed {
        uint64 Id PK
        int ChickenCategory
        int FeedType
        float64 TotalFeed
    }
    CageFeedDetail {
        uint64 Id PK
        uint64 CageFeedId FK
        uint64 ItemId FK
        float64 Percentage
    }
    CageFeedStock {
        uint64 Id PK
        uint64 CageId FK
        float64 TotalFeed
        float64 UsedFeed
    }

    CageFeed ||--o{ CageFeedDetail : "ingredients"
    Item     ||--o{ CageFeedDetail : "uses"
    Cage     ||--o{ CageFeedStock  : "has stock"

    %% ============================================================
    %% Egg monitoring
    %% ============================================================
    EggMonitoring {
        uint64 Id PK
        uint64 ChickenCageId FK
        uint64 WarehouseId FK
        uint64 TotalCrackedEgg
        uint64 TotalGoodEgg
        uint64 TotalRejectEgg
        float64 TotalWeightGoodEgg
        float64 TotalWeightCrackedEgg
    }

    ChickenCage ||--o{ EggMonitoring : "produces"
    Warehouse   ||--o{ EggMonitoring : "stored at"

    %% ============================================================
    %% Items, pricing, suppliers
    %% ============================================================
    Item {
        uint64 Id PK
        string Name
        int Category
        string Unit
        float64 DailySpending
    }
    ItemPrice {
        uint64 Id PK
        uint64 ItemId FK
        string Category
        decimal Price
        int SaleUnit
    }
    ItemPriceDiscount {
        uint64 Id PK
        string Name
        uint64 MinimumTransactionUser
        float64 TotalDiscount
    }
    Supplier {
        uint64 Id PK
        string Name
        string PhoneNumber
        string Address
        int SupplierType
    }
    SupplierItem {
        uint64 Id PK
        uint64 SupplierId FK
        uint64 ItemId FK
    }

    Item     ||--o{ ItemPrice    : "priced by"
    Supplier ||--o{ SupplierItem : "supplies"
    Item     ||--o{ SupplierItem : "supplied by"

    %% ============================================================
    %% Warehouse + warehouse items
    %% ============================================================
    Warehouse {
        uint64 Id PK
        uint64 LocationId FK
        string Name
        float64 CornCapacity
    }
    WarehouseItem {
        uint64 ItemId PK,FK
        uint64 WarehouseId PK,FK
        float64 Quantity
        time ExpiredAt
    }
    WarehouseItemHistory {
        uint64 Id PK
        string ItemName
        string ItemUnit
        string Source
        string Destination
        float64 QuantityBefore
        float64 QuantityAfter
        int Status
        uuid UserId FK
    }

    Location  ||--o{ Warehouse            : "hosts"
    Warehouse ||--o{ WarehouseItem        : "stocks"
    Item      ||--o{ WarehouseItem        : "stocked as"
    User      ||--o{ WarehouseItemHistory : "actor"

    %% ============================================================
    %% Warehouse item procurement (feed materials etc.)
    %% ============================================================
    WarehouseItemProcurementDraft {
        uint64 Id PK
        uint64 WarehouseId FK
        uint64 ItemId FK
        int64 SupplierId FK "nullable"
        float64 DailySpending
        uint64 DaysNeed
        decimal Price
    }
    WarehouseItemProcurement {
        uint64 Id PK
        uint64 WarehouseId FK
        uint64 ItemId FK
        uint64 SupplierId FK
        float64 Quantity
        float64 ReceiveQuantity
        decimal Price
        decimal TotalPrice
        time EstimationArrivalDate
        bool IsArrived
        int Status
        int PaymentStatus
        time DeadlinePaymentDate
        time PaidDate
        uuid CreatedBy FK
    }
    WarehouseItemProcurementPayment {
        uint64 Id PK
        uint64 WarehouseItemProcurementId FK
        decimal Nominal
        time PaymentDate
        uuid CreatedBy FK
    }

    Warehouse ||--o{ WarehouseItemProcurementDraft   : "drafts"
    Item      ||--o{ WarehouseItemProcurementDraft   : "for item"
    Supplier  ||--o{ WarehouseItemProcurementDraft   : "from"
    Warehouse ||--o{ WarehouseItemProcurement        : "procures into"
    Item      ||--o{ WarehouseItemProcurement        : "for item"
    Supplier  ||--o{ WarehouseItemProcurement        : "supplies"
    User      ||--o{ WarehouseItemProcurement        : "created"
    WarehouseItemProcurement ||--o{ WarehouseItemProcurementPayment : "paid by"
    User      ||--o{ WarehouseItemProcurementPayment : "created"

    %% ============================================================
    %% Corn (separate procurement flow)
    %% ============================================================
    WarehouseItemCorn {
        uint64 Id PK
        uint64 WarehouseId FK
        uint64 SupplierId FK
        float64 Quantity
        time OrderDate
        time ExpiredAt
    }
    WarehouseItemCornPrice {
        uint64 Id PK
        float64 UpperLimit
        float64 BottomLimit
        decimal BasePrice
        float64 Discount
    }
    WarehouseItemCornProcurementDraft {
        uint64 Id PK
        uint64 WarehouseId FK
        int64 SupplierId FK "nullable"
        int OvenCondition
        float64 CornWaterLevel
        bool IsOvenCanOperateInNearDay
        float64 Quantity
        decimal Price
        float64 Discount
    }
    WarehouseItemCornProcurement {
        uint64 Id PK
        uint64 WarehouseId FK
        uint64 SupplierId FK
        float64 Quantity
        float64 ReceiveQuantity
        decimal Price
        decimal TotalPrice
        int OvenCondition
        float64 CornWaterLevel
        time ExpiredAt
        time DeadlinePaymentDate
        time PaidDate
        int Status
        int PaymentStatus
        float64 Discount
        uuid CreatedBy FK
    }
    WarehouseItemCornProcurementPayment {
        uint64 Id PK
        uint64 WarehouseItemCornProcurementId FK
        decimal Nominal
        time PaymentDate
        uuid CreatedBy FK
    }

    Warehouse ||--o{ WarehouseItemCorn                    : "stocks corn"
    Supplier  ||--o{ WarehouseItemCorn                    : "supplies"
    Warehouse ||--o{ WarehouseItemCornProcurementDraft    : "drafts"
    Supplier  ||--o{ WarehouseItemCornProcurementDraft    : "from"
    Warehouse ||--o{ WarehouseItemCornProcurement         : "procures corn"
    Supplier  ||--o{ WarehouseItemCornProcurement         : "supplies"
    User      ||--o{ WarehouseItemCornProcurement         : "created"
    WarehouseItemCornProcurement ||--o{ WarehouseItemCornProcurementPayment : "paid by"
    User      ||--o{ WarehouseItemCornProcurementPayment  : "created"

    %% ============================================================
    %% Store + store items
    %% ============================================================
    Store {
        uint64 Id PK
        uint64 LocationId FK
        string Name
    }
    StoreItem {
        uint64 Id PK
        uint64 StoreId FK
        uint64 ItemId FK
        float64 Quantity
    }
    StoreItemHistory {
        uint64 Id PK
        string ItemName
        string ItemUnit
        string Source
        string Destination
        float64 QuantityBefore
        float64 QuantityAfter
        int Status
        uuid UserId FK
    }
    StoreRequestItem {
        uint64 Id PK
        uint64 WarehouseId FK
        uint64 ItemId FK
        uint64 StoreId FK
        uuid CreatedBy FK
    }

    Location  ||--o{ Store            : "hosts"
    Store     ||--o{ StoreItem        : "stocks"
    Item      ||--o{ StoreItem        : "stocked as"
    User      ||--o{ StoreItemHistory : "actor"
    Warehouse ||--o{ StoreRequestItem : "fulfills"
    Item      ||--o{ StoreRequestItem : "requested item"
    Store     ||--o{ StoreRequestItem : "requested by"
    User      ||--o{ StoreRequestItem : "created"

    %% ============================================================
    %% Customers + sales (store + warehouse channels)
    %% ============================================================
    Customer {
        uint64 Id PK
        string Name
        string PhoneNumber
    }
    StoreSale {
        uint64 Id PK
        int64 CustomerId FK "nullable"
        uint64 ItemId FK
        uint64 StoreId FK
        float64 Quantity
        decimal TotalPrice
        int PaymentStatus
        uuid CreatedBy FK
    }
    StoreSalePayment {
        uint64 Id PK
        uint64 StoreSaleId FK
        decimal Nominal
        time PaymentDate
        uuid CreatedBy FK
    }
    StoreSaleQueue {
        uint64 Id PK
        int64 CustomerId FK "nullable"
        uint64 ItemId FK
        uint64 StoreId FK
        int SaleUnit
        float64 Quantity
    }
    WarehouseSale {
        uint64 Id PK
        int64 CustomerId FK "nullable"
        uint64 ItemId FK
        uint64 WarehouseId FK
        float64 Quantity
        decimal TotalPrice
        int PaymentStatus
        uuid CreatedBy FK
    }
    WarehouseSalePayment {
        uint64 Id PK
        uint64 WarehouseSaleId FK
        decimal Nominal
        time PaymentDate
        uuid CreatedBy FK
    }
    WarehouseSaleQueue {
        uint64 Id PK
        int64 CustomerId FK "nullable"
        uint64 ItemId FK
        uint64 WarehouseId FK
        int SaleUnit
        float64 Quantity
    }

    Customer  ||--o{ StoreSale          : "buys"
    Item      ||--o{ StoreSale          : "of item"
    Store     ||--o{ StoreSale          : "at store"
    User      ||--o{ StoreSale          : "created"
    StoreSale ||--o{ StoreSalePayment   : "paid by"
    User      ||--o{ StoreSalePayment   : "created"
    Customer  ||--o{ StoreSaleQueue     : "queued for"
    Item      ||--o{ StoreSaleQueue     : "of item"
    Store     ||--o{ StoreSaleQueue     : "at store"

    Customer       ||--o{ WarehouseSale        : "buys"
    Item           ||--o{ WarehouseSale        : "of item"
    Warehouse      ||--o{ WarehouseSale        : "at warehouse"
    User           ||--o{ WarehouseSale        : "created"
    WarehouseSale  ||--o{ WarehouseSalePayment : "paid by"
    User           ||--o{ WarehouseSalePayment : "created"
    Customer       ||--o{ WarehouseSaleQueue   : "queued for"
    Item           ||--o{ WarehouseSaleQueue   : "of item"
    Warehouse      ||--o{ WarehouseSaleQueue   : "at warehouse"

    %% ============================================================
    %% Afkir (cull) chicken sales
    %% ============================================================
    AfkirChickenCustomer {
        uint64 Id PK
        string Name
        string PhoneNumber
        string Address
        decimal LatestPrice
    }
    AfkirChickenSaleDraft {
        uint64 Id PK
        uint64 AfkirChickenCustomerId FK
        uint64 ChickenCageId FK
        uint64 TotalSellChicken
        decimal PricePerChicken
        decimal TotalPrice
    }
    AfkirChickenSale {
        uint64 Id PK
        uint64 AfkirChickenCustomerId FK
        uint64 ChickenCageId FK
        uint64 TotalSellChicken
        decimal PricePerChicken
        decimal TotalPrice
        int PaymentStatus
        uuid CreatedBy FK
    }
    AfkirChickenSalePayment {
        uint64 Id PK
        uint64 AfkirChickenSaleId FK
        decimal Nominal
        time PaymentDate
        uuid CreatedBy FK
    }

    AfkirChickenCustomer ||--o{ AfkirChickenSaleDraft     : "buys"
    ChickenCage          ||--o{ AfkirChickenSaleDraft     : "sells from"
    AfkirChickenCustomer ||--o{ AfkirChickenSale          : "buys"
    ChickenCage          ||--o{ AfkirChickenSale          : "sells from"
    User                 ||--o{ AfkirChickenSale          : "created"
    AfkirChickenSale     ||--o{ AfkirChickenSalePayment   : "paid by"
    User                 ||--o{ AfkirChickenSalePayment   : "created"

    %% ============================================================
    %% Work (daily + additional) + attendance + salary + cash advance
    %% ============================================================
    DailyWork {
        uint64 Id PK
        string Description
        uint64 RoleId FK
        time StartTime
        time EndTime
    }
    DailyWorkUser {
        uint64 Id PK
        uint64 DailyWorkId FK
        uuid UserId FK
        bool IsDone
    }
    AdditionalWork {
        uint64 Id PK
        string Name
        uint64 LocationId FK
        int64 WarehouseId FK "nullable"
        int64 StoreId FK "nullable"
        int64 CageId FK "nullable"
        string Description
        uint64 Slot
        time WorkDate
        decimal Salary
    }
    AdditionalWorkUser {
        uint64 Id PK
        uuid UserId FK
        uint64 AdditionalWorkId FK
        bool IsDone
    }
    UserPresence {
        uint64 Id PK
        uuid UserId FK
        time StartTime
        time EndTime
        int Status
        string Note
        string Evidence
        int SubmissionPresenceStatus
    }
    UserSalaryPayment {
        uint64 Id PK
        uuid UserId FK
        decimal BaseSalary
        decimal BonusSalary
        decimal CompentationSalary
        decimal AdditionalWorkSalary
        decimal Cashbond
        string PaymentProof
        int PaymentMethod
        time PaymentDate
        bool IsPaid
        uuid CreatedBy FK
    }
    UserCashAdvance {
        uint64 Id PK
        uuid UserId FK
        decimal Nominal
        time DeadlinePaymentDate
        int PaymentStatus
        time PaidDate
        uuid CreatedBy FK
    }
    UserCashAdvancePayment {
        uint64 Id PK
        uint64 UserCashAdvanceId FK
        decimal Nominal
        time PaymentDate
        uuid CreatedBy FK
    }

    Role             ||--o{ DailyWork              : "scoped to"
    DailyWork        ||--o{ DailyWorkUser          : "assigned to"
    User             ||--o{ DailyWorkUser          : "performs"
    Location         ||--o{ AdditionalWork         : "scoped to"
    Warehouse        ||--o{ AdditionalWork         : "at warehouse"
    Store            ||--o{ AdditionalWork         : "at store"
    Cage             ||--o{ AdditionalWork         : "at cage"
    AdditionalWork   ||--o{ AdditionalWorkUser     : "assigned to"
    User             ||--o{ AdditionalWorkUser     : "performs"
    User             ||--o{ UserPresence           : "attends"
    User             ||--o{ UserSalaryPayment      : "earns"
    User             ||--o{ UserCashAdvance        : "borrows"
    UserCashAdvance  ||--o{ UserCashAdvancePayment : "repaid by"
    User             ||--o{ UserCashAdvancePayment : "created"
    User             ||--o{ UserSalaryPayment      : "created"
    User             ||--o{ UserCashAdvance        : "created"

    %% ============================================================
    %% Misc: expenses, cashflow snapshot, notifications
    %% ============================================================
    Expense {
        uint64 Id PK
        int ExpenseCategory
        string Name
        string ReceiverName
        string ReceiverPhoneNumber
        decimal Nominal
        int PaymentMethod
        string PaymentProof
        string Description
        uint64 LocationId FK
        int64 WarehouseId FK "nullable"
        int64 StoreId FK "nullable"
        int64 CageId FK "nullable"
        int LocationType
        uuid CreatedBy FK
    }
    CashflowHistory {
        uint64 Id PK
        uint64 LocationId FK
        decimal Income
        decimal Profit
        decimal Expense
        decimal Cash
        decimal Receivables
        decimal Debt
        decimal StoreEggSale
        decimal WarehouseEggSale
    }
    Notification {
        uint64 Id PK
        uuid UserId FK "nullable"
        int64 StoreId FK "nullable"
        int64 WarehouseId FK "nullable"
        int64 CageId FK "nullable"
        string Description
        text_array NotificationContexts
        bool IsMarked
    }

    Location  ||--o{ Expense         : "scoped to"
    Warehouse ||--o{ Expense         : "at warehouse"
    Store     ||--o{ Expense         : "at store"
    Cage      ||--o{ Expense         : "at cage"
    User      ||--o{ Expense         : "created"
    Location  ||--o{ CashflowHistory : "snapshot of"
    User      ||--o{ Notification    : "notified"
    Store     ||--o{ Notification    : "about store"
    Warehouse ||--o{ Notification    : "about warehouse"
    Cage      ||--o{ Notification    : "about cage"
```

## Notes

- **`CreatedBy` / `UpdatedBy` audit columns** exist on most tables as `uuid.NullUUID` referencing `User.Id`. Only `CreatedBy` edges are drawn (and only where a `CreatedByUser` GORM relation is declared) to keep the diagram readable. Mentally, every entity also has an "updated by User" edge.
- **Nullable FKs** (those backed by `sql.NullInt64` / `uuid.NullUUID`) are annotated `"nullable"`. Notification and Expense in particular have polymorphic-style optional FKs to `Cage`, `Store`, `Warehouse`, with a `LocationType` discriminator on `Expense`.
- **Soft delete** (`gorm.DeletedAt`) is applied to `Cage`, `ChickenCage`, `AdditionalWork`, `DailyWork`.
- **No FK at all**: `ChickenHealthItem` (lookup), `ItemPriceDiscount` (lookup), `WarehouseItemCornPrice` (lookup), `CageFeedStock.CageId` and `ChickenPerformance.LocationId` are referenced by ID but have no GORM `foreignKey` relation declared — the diagram still shows the logical edge.
- **Composite-key tables** (`CagePlacement`, `StorePlacement`, `WarehousePlacement`, `WarehouseItem`) use `(UserId, CageId)` / `(UserId, StoreId)` / `(ItemId, WarehouseId)` as their primary key.
- **History tables** (`WarehouseItemHistory`, `StoreItemHistory`) intentionally store denormalized item name/unit instead of an `ItemId` FK — the source/destination strings document this design choice.
- **Drafts** (`ChickenProcurementDraft`, `WarehouseItemProcurementDraft`, `WarehouseItemCornProcurementDraft`, `AfkirChickenSaleDraft`) are scratchpads that are promoted to their non-draft counterpart on confirmation.
