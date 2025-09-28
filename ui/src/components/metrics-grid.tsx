import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  TrendingUp,
  TrendingDown,
  Activity,
  Zap,
  Thermometer,
  Wifi,
} from "lucide-react";

const metrics = [
  {
    title: "Active Devices",
    value: "24",
    change: "+2",
    trend: "up",
    icon: Activity,
    color: "text-chart-1",
  },
  {
    title: "Power Consumption",
    value: "1.2kW",
    change: "-5%",
    trend: "down",
    icon: Zap,
    color: "text-chart-3",
  },
  {
    title: "Avg Temperature",
    value: "22.5°C",
    change: "+1.2°C",
    trend: "up",
    icon: Thermometer,
    color: "text-chart-2",
  },
  {
    title: "Network Status",
    value: "99.9%",
    change: "Stable",
    trend: "stable",
    icon: Wifi,
    color: "text-chart-4",
  },
];

export default function MetricsGrid() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {metrics.map((metric, index) => (
        <Card key={index} className="bg-card border-border">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {metric.title}
            </CardTitle>
            <metric.icon className={`w-4 h-4 ${metric.color}`} />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metric.value}</div>
            <div className="flex items-center space-x-2 text-xs text-muted-foreground">
              {metric.trend === "up" && (
                <TrendingUp className="w-3 h-3 text-chart-2" />
              )}
              {metric.trend === "down" && (
                <TrendingDown className="w-3 h-3 text-chart-1" />
              )}
              <span
                className={
                  metric.trend === "up"
                    ? "text-chart-2"
                    : metric.trend === "down"
                    ? "text-chart-1"
                    : "text-muted-foreground"
                }
              >
                {metric.change}
              </span>
              <span>from last hour</span>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
