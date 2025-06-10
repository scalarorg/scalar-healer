-- name: GetUtxoSnapshot :many
SELECT
    u.tx_id,
    u.vout,
    u.script_pubkey,
    u.amount_in_sats,
    u.custodian_group_uid,
    u.block_height,
    COALESCE(
      json_agg(
        CASE
          WHEN r.request_id IS NOT NULL THEN
            json_build_object('request_id', r.request_id, 'amount', ur.amount)
        END
      ) FILTER (WHERE r.request_id IS NOT NULL),
      '[]'
    ) AS reservations
  FROM utxos u
  LEFT JOIN utxo_reservations ur ON u.tx_id = ur.utxo_tx_id AND u.vout = ur.utxo_vout
  LEFT JOIN reservations r ON r.request_id = ur.reservation_id
  WHERE u.custodian_group_uid = $1::bytea
  GROUP BY u.tx_id, u.vout, u.script_pubkey, u.amount_in_sats, u.custodian_group_uid, u.block_height
  ORDER BY u.tx_id;
