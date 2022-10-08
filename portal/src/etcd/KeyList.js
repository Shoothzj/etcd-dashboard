import React, { useEffect, useState } from 'react'
import { DataGrid } from '@mui/x-data-grid';
import { Box } from '@mui/material';
import { GridEventListener } from '@mui/x-data-grid';
import { useNavigate } from 'react-router-dom';
import BACKEND_HOST from '../const';


const KeyListFetch = () => {

  const columns = [
    { field: 'key', headerName: 'Key', width: 300 },
  ]

  const [keys, setKeys] = useState([])

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
    navigate(`/etcd/keys/${params.row.key}`)
  };

  return (
    <div>
      <h1>Keys</h1>
      <Box sx={{ height: 400, width: '100%' }}>
        <DataGrid onRowClick={handleEvent}
          rows={keys}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10]}
          checkboxSelection
          disableSelectionOnClick
          experimentalFeatures={{ newEditingApi: true }}
        />
      </Box>
    </div>
  )
}

export default KeyListFetch
