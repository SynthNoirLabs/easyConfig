import type { LucideIcon } from "lucide-react";
import React from "react";
import "../styles/empty-state.css";

interface EmptyStateProps {
  icon: LucideIcon;
  title: string;
  description: string;
  action?: {
    label: string;
    onClick: () => void;
  };
}

export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
}: EmptyStateProps) {
  return (
    <div className="empty-state">
      <Icon size={48} />
      <h3>{title}</h3>
      <p>{description}</p>
      {action && <button type="button" onClick={action.onClick}>{action.label}</button>}
    </div>
  );
}
