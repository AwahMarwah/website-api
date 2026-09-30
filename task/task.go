package task

// TypeReconcilePayments menyelaraskan status order dengan status sebenarnya di
// payment provider, untuk order yang webhook-nya tidak pernah sampai.
const TypeReconcilePayments = "order:reconcile_payments"
