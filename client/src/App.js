import React, { useState } from "react";
import {
  Container,
  Typography,
  TextField,
  Button,
  Box,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from "@mui/material";
import { Delete } from "@mui/icons-material";

function App() {
  // Default pack sizes
  const [packs, setPacks] = useState([250, 500, 1000, 2000, 5000]);
  // Default items
  const [items, setItems] = useState(12001);
  // Result from server (pack -> quantity)
  const [result, setResult] = useState(null);
  // Error handling
  const [error, setError] = useState("");

  // Modal state for adding new pack size
  const [openModal, setOpenModal] = useState(false);
  const [newPackValue, setNewPackValue] = useState("");

  // Handle pack size change and clear previous result
  const handlePackChange = (index, value) => {
    // Clear previous calculation result
    setResult(null);

    const numValue = parseInt(value, 10);
    if (isNaN(numValue) || numValue <= 0) {
      const newPacks = [...packs];
      newPacks[index] = "";
      setPacks(newPacks);
      return;
    }
    const newPacks = [...packs];
    newPacks[index] = numValue;
    setPacks(newPacks);
  };

  // Delete a pack size and clear previous result
  const deletePackSize = (index) => {
    setResult(null);
    const newPacks = packs.filter((_, i) => i !== index);
    setPacks(newPacks);
  };

  // Handle items change and clear previous result
  const handleItemsChange = (value) => {
    setResult(null);
    const numValue = parseInt(value, 10);
    if (isNaN(numValue) || numValue <= 0) {
      setItems("");
      return;
    }
    setItems(numValue);
  };

  // Calculate packs by sending POST to localhost:5000
  const handleCalculate = async () => {
    // Clear previous error and result
    setError("");
    setResult(null);

    if (!packs.length || !items) {
      setError("Please provide valid pack sizes and item count.");
      return;
    }

    try {
      const response = await fetch("http://localhost:5000", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ packs, order: items }),
      });
      if (!response.ok) {
        const text = await response.text();
        throw new Error(text);
      }
      const data = await response.json();
      setResult(data);
    } catch (err) {
      setError(err.message);
    }
  };

  // Open the modal for adding a new pack size
  const handleOpenModal = () => {
    setNewPackValue("");
    setOpenModal(true);
  };

  // Close the modal and move focus to document body to prevent focus errors
  const handleCloseModal = () => {
    setOpenModal(false);
    document.body.focus();
  };

  // Handle modal input change
  const handleNewPackChange = (e) => {
    setNewPackValue(e.target.value);
  };

  // Add new pack size from modal if valid, then clear previous result
  const handleAddNewPack = () => {
    const numValue = parseInt(newPackValue, 10);
    if (isNaN(numValue) || numValue <= 0) {
      handleCloseModal();
      return;
    }
    setPacks([...packs, numValue]);
    setResult(null);
    handleCloseModal();
  };

  // Determine if the new pack value is valid (positive number)
  const isValidNewPack = newPackValue && parseInt(newPackValue, 10) > 0;

  return (
    <Container maxWidth="sm" sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Order Packs Calculator
      </Typography>

      {/* Pack Sizes Section */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h6">Pack Sizes</Typography>
        {packs.map((pack, index) => (
          <Box
            key={index}
            sx={{ display: "flex", alignItems: "center", mt: 1 }}
          >
            <TextField
              label={`Pack #${index + 1}`}
              variant="outlined"
              size="small"
              type="number"
              value={pack}
              onChange={(e) => handlePackChange(index, e.target.value)}
              sx={{ mr: 1 }}
            />
            <IconButton color="error" onClick={() => deletePackSize(index)}>
              <Delete />
            </IconButton>
          </Box>
        ))}
        <Button
          variant="contained"
          color="primary"
          onClick={handleOpenModal}
          sx={{ mt: 2 }}
        >
          Add Pack Size
        </Button>
      </Box>

      {/* Modal for Adding Pack Size */}
      <Dialog
        open={openModal}
        onClose={handleCloseModal}
        closeAfterTransition={true}
      >
        <DialogTitle>Enter Pack Size</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Pack Size"
            type="number"
            fullWidth
            variant="outlined"
            value={newPackValue}
            onChange={handleNewPackChange}
            slotProps={{ htmlInput: { min: 1 } }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseModal} variant="outlined">
            Cancel
          </Button>
          <Button
            onClick={handleAddNewPack}
            variant="outlined"
            disabled={!isValidNewPack}
          >
            Add
          </Button>
        </DialogActions>
      </Dialog>

      {/* Items Section */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h6">Calculate packs for order</Typography>
        <TextField
          label="Items"
          variant="outlined"
          type="number"
          value={items}
          onChange={(e) => handleItemsChange(e.target.value)}
          size="small"
          sx={{ mr: 2, mt: 1 }}
        />
        <Button
          variant="contained"
          color="success"
          onClick={handleCalculate}
          sx={{ mt: 1 }}
        >
          Calculate
        </Button>
      </Box>

      {/* Error Message */}
      {error && (
        <Typography variant="body1" color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}

      {/* Result Table */}
      {result && (
        <Box sx={{ mt: 3 }}>
          <Typography variant="h6">Result:</Typography>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>
                  <strong>Pack</strong>
                </TableCell>
                <TableCell>
                  <strong>Quantity</strong>
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {Object.entries(result).map(([pack, quantity]) => (
                <TableRow key={pack}>
                  <TableCell>{pack}</TableCell>
                  <TableCell>{quantity}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Box>
      )}
    </Container>
  );
}

export default App;
