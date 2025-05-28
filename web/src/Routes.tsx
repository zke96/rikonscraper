import { Box } from "@mui/material";
import { Route, Routes } from "react-router";
import AlertsPage from "./AlertsPage";
import Home from "./Home";

export function AppRoutes() {
    return (
        <Box sx={{ position: 'relative', top: '64px', padding: '1em', display: 'inline-block' }}>
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/alerts/:email" element={<AlertsPage />} />
            </Routes>
        </Box>
    )
}