import { AppBar, Autocomplete, Box, Card, CardContent, CardHeader, Link, Skeleton, Stack, TextField, Toolbar, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import './Home.css';

interface Product {
    label: string;
    url: string;
    price?: string;
    inStock?: boolean;
    description?: string;
}

function Home() {
    const [products, setProducts] = useState<Product[]>([])
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
    const [parts, setParts] = useState<Product[]>([])
    const [selectedPart, setSelectedPart] = useState<Product | null>(null)
    const [productDetails, setProductDetails] = useState<Product | null>(null)
    const [detailsLoading, setDetailsLoading] = useState<boolean>(false);

    useEffect(() => {
        fetch("http://localhost:8080/products").then(response => {
            response.json().then((json: Product[]) => {
                setProducts(json)
            })
        })
    }, [])

    useEffect(() => {
        if (selectedProduct) {
            setParts([]);
            setSelectedPart(null);
            fetch(`http://localhost:8080/parts?partNumber=${selectedProduct.url}`).then(response => {
                response.json().then((json: Product[]) => {
                    setParts(json)
                })
            })
        }
    }, [selectedProduct])

    useEffect(() => {
        if (selectedPart) {
            setDetailsLoading(true);
            setProductDetails(null);
            fetchPartInfo(selectedPart)
        }
    }, [selectedPart])

    const fetchPartInfo = (p: Product) => {
        fetch(`http://localhost:8080/partinfo?url=${p.url}`).then(response => {
            response.json().then(json => {
                console.log(json)
                setProductDetails(json)
                setDetailsLoading(false)
            })
        }).catch(() => setDetailsLoading(false))
    }

    return (
        <Box sx={{ background: 'background.paper' }}>
            <AppBar>
                <Toolbar>
                    <Typography
                        variant="h6"
                        noWrap
                    >
                        Rikon Parts Stock Alert
                    </Typography>
                </Toolbar>
            </AppBar>
            <Stack spacing={2} sx={{ alignItems: 'center', width: '400px' }}>
                <Autocomplete
                    fullWidth
                    options={products}
                    value={selectedProduct}
                    loading={products.length === 0}
                    loadingText="Loading Products..."
                    onChange={(_e, v) => v && setSelectedProduct(v)}
                    renderInput={(params) => <TextField {...params} label="Select Product" />}
                />
                {selectedProduct &&
                    <Autocomplete
                        fullWidth
                        options={parts}
                        value={selectedPart}
                        loading={parts.length === 0}
                        loadingText="Loading Parts..."
                        onChange={(_e, v) => setSelectedPart(v)}
                        renderInput={(params) => <TextField {...params} label="Select Part" />}
                    />}
                {productDetails && <Card sx={{ width: '500px' }}>
                    <CardHeader
                        classes={{
                            title: 'title',
                            subheader: 'title'
                        }}
                        title={productDetails.label} subheader={productDetails.description} sx={{ borderBottom: '1px solid', borderColor: 'divider' }} />
                    <CardContent>
                        <Stack sx={{ alignItems: 'flex-start' }}>
                            <Typography>Price: ${productDetails.price}</Typography>
                            <Typography>Stock Status: {productDetails.inStock ? 'In Stock' : 'Out of Stock'}</Typography>
                            {productDetails.url && <Link href={`https://rikontools.com/products/${productDetails.url}`}>View Part Page</Link>}
                        </Stack>
                    </CardContent>
                </Card>}
                {detailsLoading && <Skeleton variant="rounded" width={300} height={120} />}
            </Stack>
        </Box>
    )
}

export default Home