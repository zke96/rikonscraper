import { Autocomplete, Button, Card, CardContent, CardHeader, Link, Skeleton, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import './Home.css';
import { getPartDetails, getPartsForProduct, getProducts, putAlert } from "./util/api";
import type { Product } from "./util/types";

function Home() {
    const [products, setProducts] = useState<Product[]>([])
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
    const [parts, setParts] = useState<Product[]>([])
    const [selectedPart, setSelectedPart] = useState<Product | null>(null)
    const [productDetails, setProductDetails] = useState<Product | null>(null)
    const [detailsLoading, setDetailsLoading] = useState<boolean>(false);
    const [email, setEmail] = useState<string>('');

    useEffect(() => {
        getProducts().then((products) => {
            console.log(products)
            setProducts(products)
        })
    }, [])

    useEffect(() => {
        if (selectedProduct) {
            setParts([]);
            setSelectedPart(null);
            getPartsForProduct(selectedProduct.id).then(parts => {
                setParts(parts)
            })
        }
    }, [selectedProduct])

    useEffect(() => {
        if (selectedPart) {
            setDetailsLoading(true);
            setProductDetails(null);
            getPartDetails(selectedPart).then(details => {
                setProductDetails(details);
                setDetailsLoading(false);
            })
        }
    }, [selectedPart])

    const submitNewAlert = () => {
        if (selectedPart) {
            putAlert(email, selectedPart.id).then(() => {
                window.alert(`Successfully subscribed to updates for product ${selectedPart.label}`)
            }).catch(() => {
                window.alert('Failed to subscribe')
            })
            setEmail('');
        }
    }

    return (
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
                        <Typography>Enter your email to be notified when this part is back in stock</Typography>
                        <TextField value={email} onChange={(e) => setEmail(e.target.value)} />
                        <Button onClick={submitNewAlert}>Submit</Button>
                    </Stack>
                </CardContent>
            </Card>}
            {detailsLoading && <Skeleton variant="rounded" width={300} height={120} />}
        </Stack>
    )
}

export default Home