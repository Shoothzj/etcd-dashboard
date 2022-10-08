import React, { useEffect, useState } from 'react'
import { DataGrid, GridEventListener, GridToolbar } from '@mui/x-data-grid';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import BACKEND_HOST from '../const';

const KeyListFetch = () => {

  const columns = [
    { field: 'key', headerName: 'Key', width: 300 },
  ]

  const [keys, setKeys] = useState([])

  const [dialogOpen, setDialogOpen] = React.useState(false);

  const [key, setKey] = React.useState('');

  const [value, setValue] = React.useState('');

  const handleKeyChanged = (event) => {
    setKey(event.target.value);
  };

  const handleValueChanged = (event) => {
    setValue(event.target.value);
  };

  const handleClickOpen = () => {
    setDialogOpen(true);
  };

  const handlePutKey = () => {
    fetch(BACKEND_HOST + '/api/etcd/keys', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        key: key,
        value: value
      })
    })
      .then(response => response.json())
      .then(data => {
        console.log('Success:', data);
        setDialogOpen(false);
        setKey('');
        setValue('');
        fetchKeys();
      })
      .catch((error) => {
        console.error('Error:', error);
      });
    setDialogOpen(false);
  };

  const navigate = useNavigate();

  const fetchKeys = async () => {
    const response = await fetch(BACKEND_HOST + '/api/etcd/keys')
    const data = await response.json()
    setKeys(data.map((key) => ({ id: key, key: key })))
  }

  useEffect(() => {
    fetchKeys()
  }, [])

  const handleEvent: GridEventListener<'rowClick'> = (
    params, // GridRowParams
    event, // MuiEvent<React.MouseEvent<HTMLElement>>
    details, // GridCallbackDetails
  ) => {
    navigate(`/etcd/keys/${encodeURIComponent(params.row.key)}`)
  };

  return (
    <div>
      <h1>Keys</h1>
      <Button variant="contained" onClick={handleClickOpen}>
        Put Key
      </Button>
      <Dialog open={dialogOpen} onClose={handlePutKey}>
        <DialogTitle>Subscribe</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Please form the key and value
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="put-key"
            label="Key"
            value={key}
            onChange={handleKeyChanged}
            type="text"
            fullWidth
            variant="standard"
          />
          <TextField
            autoFocus
            margin="dense"
            id="put-value"
            label="Value"
            value={value}
            onChange={handleValueChanged}
            type="text"
            fullWidth
            variant="standard"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handlePutKey}>Cancel</Button>
          <Button onClick={handlePutKey}>Confirm</Button>
        </DialogActions>
      </Dialog>
      <Box sx={{ height: 600, width: '100%' }}>
        <DataGrid onRowClick={handleEvent}
                  rows={keys}
                  columns={columns}
                  pageSize={10}
                  rowsPerPageOptions={[10]}
                  checkboxSelection
                  disableSelectionOnClick
                  experimentalFeatures={{ newEditingApi: true }}
                  components={{ Toolbar: GridToolbar }}
        />
      </Box>
    </div>
  )
}

export default KeyListFetch
