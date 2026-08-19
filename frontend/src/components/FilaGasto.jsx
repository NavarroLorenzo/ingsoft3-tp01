import { formatearFecha, formatearMoneda } from "../utils/formato";

export default function FilaGasto({ gasto, onEditar, onEliminar }) {
  return (
    <tr>
      <td>{gasto.descripcion}</td>
      <td>{gasto.categoria?.nombre}</td>
      <td>{formatearFecha(gasto.fecha)}</td>
      <td className="amount">{formatearMoneda(gasto.monto)}</td>
      <td className="row-actions">
        <button className="text-button" type="button" onClick={() => onEditar(gasto)}>Editar</button>
        <button className="text-button danger" type="button" onClick={() => onEliminar(gasto)}>Eliminar</button>
      </td>
    </tr>
  );
}
