import React from 'react';
import { type EnrichedContent } from '../api/apiSlice';
import {
    Box,
    List,
    ListItem,
    ListItemText,
    Alert
} from '@mui/material';

interface PlantingDisplayProps {
    contents: EnrichedContent[];
}

const PlantingDisplay: React.FC<PlantingDisplayProps> = ({ contents }) => {
    if (contents.length === 0) {
        return <Alert severity="info">На этой грядке пока ничего не растет.</Alert>;
    }

    return (
        <Box>
            <List dense>
                {contents.map((content) => (
                    <ListItem key={content.item.id}>
                        <ListItemText
                            primary={`${content.item.name} (x${content.quantity})`}
                            secondary={content.item.description}
                        />
                    </ListItem>
                ))}
            </List>
        </Box>
    );
};

export default PlantingDisplay;