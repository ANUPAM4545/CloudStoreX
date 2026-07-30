import React from "react";
import { Card, CardContent } from "@/components/ui/card";

export interface EmptyStateProps {
  icon?: React.ReactNode;
  title: string;
  description: string;
  action?: React.ReactNode;
}

export function EmptyState({ icon, title, description, action }: EmptyStateProps) {
  return (
    <Card className="border-dashed shadow-sm bg-background/50">
      <CardContent className="flex flex-col items-center justify-center p-12 text-center min-h-[300px]">
        {icon && (
          <div className="h-14 w-14 rounded-full bg-muted flex items-center justify-center mb-4 text-muted-foreground">
            {icon}
          </div>
        )}
        <h3 className="text-lg font-semibold tracking-tight mb-1 text-foreground">{title}</h3>
        <p className="text-sm text-muted-foreground max-w-sm mb-6">{description}</p>
        {action && <div>{action}</div>}
      </CardContent>
    </Card>
  );
}
