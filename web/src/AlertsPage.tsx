import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import { Box, CircularProgress, IconButton, Stack, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import './Home.css';
import { getAlertsByEmail } from "./util/api";
import type { Alert } from "./util/types";

function AlertsPage() {
    const [alerts, setAlerts] = useState<Alert[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const params = useParams()
    const { email } = params

    useEffect(() => {
        if (email) {
            setLoading(true)
            getAlertsByEmail(email).then(alerts => {
                setAlerts(alerts)
                setLoading(false)
            })
        }
    }, [email])

    return (
        <Box sx={{ background: 'background.paper' }}>
            <Stack spacing={2} sx={{ alignItems: 'left', width: '600px' }}>
                <Typography variant="subtitle2" sx={{ textAlign: 'left', color: 'text.primary' }}>{`Email ${email} is currently subscribed to alerts for the following products`}</Typography>
                {alerts.length > 0 &&
                    <Stack spacing={1}>
                        {alerts.map((a) => (
                            <Box key={`alert-${a.productCode}`}>
                                <Typography>{a.productCode}</Typography>
                                <IconButton><DeleteOutlineIcon /></IconButton>
                            </Box>
                        ))}
                    </Stack>}
                {loading && <CircularProgress />}
            </Stack>
        </Box>
    )
}

export default AlertsPage