"use client";

import { Activity, Thermometer, TrendingUp } from "lucide-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { fetchMetricsData, type IMetricsData } from "@/app/page";
import { useEffect } from "react";

export default function IoTMetricsDashboard({ data }: { data: IMetricsData }) {
  const sensors = Object.entries(data.sensor_average);

  const getStatusColor = (temp: number) => {
    if (temp < 20) return "bg-blue-500";
    if (temp < 25) return "bg-green-500";
    if (temp < 30) return "bg-yellow-500";
    return "bg-red-500";
  };

  const getStatusText = (temp: number) => {
    if (temp < 20) return "Cool";
    if (temp < 25) return "Normal";
    if (temp < 30) return "Warm";
    return "Hot";
  };

  useEffect(() => {
    const fetchData = async () => {
      await fetchMetricsData();
    };
    fetchData();
  }, []);

  return (
    <div className="space-y-6">
      <Card className="border-2 border-primary/20">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Overall Average</CardTitle>
          <TrendingUp className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{data.overall_average}°C</div>
          <p className="text-xs text-muted-foreground">
            Across all {sensors.length} sensors
          </p>
        </CardContent>
      </Card>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {sensors.map(([sensorId, temperature]) => (
          <Card key={sensorId} className="relative overflow-hidden">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {sensorId.replace("sensor", "Sensor ")}
              </CardTitle>
              <div className="flex items-center gap-2">
                <Badge
                  variant="secondary"
                  className={`${getStatusColor(temperature)} text-white`}
                >
                  {getStatusText(temperature)}
                </Badge>
                <Thermometer className="h-4 w-4 text-muted-foreground" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{temperature}°C</div>
              <div className="flex items-center gap-2 mt-2">
                <Activity className="h-3 w-3 text-green-500" />
                <p className="text-xs text-muted-foreground">Active</p>
              </div>

              {/* Temperature bar indicator */}
              <div className="mt-3 w-full bg-secondary rounded-full h-2">
                <div
                  className={`h-2 rounded-full ${getStatusColor(temperature)}`}
                  style={{
                    width: `${Math.min((temperature / 40) * 100, 100)}%`,
                  }}
                />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Highest Reading
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold">
              {Math.max(...Object.values(data.sensor_average))}°C
            </div>
            <p className="text-xs text-muted-foreground">
              Sensor{" "}
              {Object.entries(data.sensor_average)
                .find(
                  ([, temp]) =>
                    temp === Math.max(...Object.values(data.sensor_average))
                )?.[0]
                .replace("sensor", "")}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Lowest Reading
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold">
              {Math.min(...Object.values(data.sensor_average))}°C
            </div>
            <p className="text-xs text-muted-foreground">
              Sensor{" "}
              {Object.entries(data.sensor_average)
                .find(
                  ([, temp]) =>
                    temp === Math.min(...Object.values(data.sensor_average))
                )?.[0]
                .replace("sensor", "")}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Temperature Range
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold">
              {(
                Math.max(...Object.values(data.sensor_average)) -
                Math.min(...Object.values(data.sensor_average))
              ).toFixed(1)}
              °C
            </div>
            <p className="text-xs text-muted-foreground">
              Variation across sensors
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
