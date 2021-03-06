package metastg

import (
	"bytes"
	"database/sql"
	"soloos/sdfs/types"

	"github.com/gocraft/dbr"
)

func (p *NetBlockDriver) FetchNetBlockFromDB(pNetINode *types.NetINode,
	netBlockIndex int, pNetBlock *types.NetBlock,
	backendPeerIDArrStr *string) (err error) {
	var (
		sess    *dbr.Session
		sqlRows *sql.Rows
	)

	sess = p.helper.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("netblock_len", "netblock_cap", "backend_peer_id_arr").
		From("b_netblock").
		Where("netinode_id=? and index_in_netinode=?",
			pNetINode.IDStr(), netBlockIndex,
		).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	pNetBlock.IndexInNetINode = netBlockIndex
	pNetBlock.NetINodeID = pNetINode.ID
	err = sqlRows.Scan(&pNetBlock.Len, &pNetBlock.Cap, backendPeerIDArrStr)
	if err != nil {
		goto QUERY_DONE
	}

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return err
}

func (p *NetBlockDriver) StoreNetBlockInDB(pNetINode *types.NetINode, pNetBlock *types.NetBlock) error {
	var (
		sess                *dbr.Session
		tx                  *dbr.Tx
		netINodeIDStr       = pNetINode.IDStr()
		backendPeerIDArrStr bytes.Buffer
		i                   int
		err                 error
	)

	sess = p.helper.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		goto QUERY_DONE
	}

	if pNetBlock.StorDataBackends.Len > 0 {
		backendPeerIDArrStr.WriteString(pNetBlock.StorDataBackends.Arr[0].Ptr().PeerIDStr())
		for i = 1; i < pNetBlock.StorDataBackends.Len; i++ {
			backendPeerIDArrStr.WriteString(",")
			backendPeerIDArrStr.WriteString(pNetBlock.StorDataBackends.Arr[i].Ptr().PeerIDStr())
		}
	} else {
		backendPeerIDArrStr.WriteString("")
	}

	_, err = sess.InsertInto("b_netblock").
		Columns("netinode_id", "index_in_netinode", "netblock_len", "netblock_cap", "backend_peer_id_arr").
		Values(netINodeIDStr, pNetBlock.IndexInNetINode, pNetBlock.Len, pNetBlock.Cap,
			backendPeerIDArrStr.String()).
		Exec()
	if err != nil {
		_, err = sess.Update("b_netblock").
			Set("netblock_len", pNetBlock.Len).
			Set("netblock_cap", pNetBlock.Cap).
			Set("backend_peer_id_arr", backendPeerIDArrStr.String()).
			Where("netinode_id=? and index_in_netinode=?", netINodeIDStr, pNetBlock.IndexInNetINode).
			Exec()
	}

QUERY_DONE:
	if err != nil {
		tx.RollbackUnlessCommitted()
	} else {
		err = tx.Commit()
	}
	return err
}
