import React from 'react';
import RaisedButton from 'material-ui/RaisedButton';
import {
  TableRow,
  TableRowColumn
} from 'material-ui/Table';

const PullRow = ({ pull, index }) => {
  const { title, body } = pull.attributes;

  return(
    <TableRow key={index} displayBorder={true}>
      <TableRowColumn>
        {title}
      </TableRowColumn>
      <TableRowColumn>
        {body}
      </TableRowColumn>
    </TableRow>
  );
};

export default PullRow;
