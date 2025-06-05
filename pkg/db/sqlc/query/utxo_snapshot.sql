-- name: GetUtxoSnapshot :many
SELECT
    u.tx_id,
    u.vout,
    u.script_pubkey,
    u.amount_in_sats,
    u.custodian_group_uid,
    json_agg(json_build_object(
      'request_id', r.request_id,
      'amount', r.amount
    )) AS reservations
  FROM utxos u
  JOIN utxo_reservations ur ON u.tx_id = ur.utxo_tx_id AND u.vout = ur.utxo_vout
  JOIN reservations r ON r.id = ur.reservation_id
  WHERE u.custodian_group_uid = $1::bytea
  GROUP BY u.tx_id, u.vout, u.script_pubkey, u.amount_in_sats, u.custodian_group_uid;