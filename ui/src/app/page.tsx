import IoTMetricsDashboard from "@/components/iot-metrics-dashboard";

import { AppSidebar } from "@/components/app-sidebar";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";

export interface IMetricsData {
  overall_average: number;
  sensor_average: Record<string, number>;
}

export async function fetchMetricsData(): Promise<IMetricsData | null> {
  try {
    const response = await fetch(`http://localhost:8080/temperature`, {
      cache: "no-store",
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch metrics: ${response.status}`);
    }

    const data = await response.json();
    return data?.data;
  } catch (error) {
    return null;
  }
}

export default async function Page() {
  const result = null;

  return (
    <SidebarProvider
      style={
        {
          "--sidebar-width": "19rem",
        } as React.CSSProperties
      }
    >
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 px-4">
          <SidebarTrigger className="-ml-1" />
          <Separator
            orientation="vertical"
            className="mr-2 data-[orientation=vertical]:h-4"
          />
        </header>

        <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
          <IoTMetricsDashboard
            data={
              result
                ? result
                : {
                    overall_average: 27.5,
                    sensor_average: {
                      sensor123: 25.5,
                      sensor456: 30.0,
                      sensor789: 27.0,
                    },
                  }
            }
          />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
