import React from 'react';
import { connect } from 'react-redux';
import { initialize } from '../actions';
import {
  Table,
  TableBody,
  TableRow,
  TableRowColumn,
  TableHeader,
  TableHeaderColumn
} from 'material-ui/Table';

import PullRow from './PullRow';

export class PullsTable extends React.Component {
  componentWillMount() {
    this.props.initialize();
  }

  headerColumnStyle() {
    return { textAlign: 'center' };
  }

  // headerColumnGenerator(text, index) {
  //   return(
  //     <TableHeaderColumn key={index} style={this.headerColumnStyle()}>
  //       {text}
  //     </TableHeaderColumn>
  //   );
  // }

  // const headerNames = ['Title', 'Body'];

  // <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
  //   <TableRow>
  //     {headerNames.map((name, index) => {
  //       return this.headerColumnGenerator(name, index);
  //     })}
  //   </TableRow>
  // </TableHeader>

  render() {
    const { pulls, isFetching } = this.props;

    return(
      <Table fixedHeader={false} selectable={false}>
        <TableBody displayRowCheckbox={false} selectable={false}>
          {isFetching ?
            <TableRow>
              <TableRowColumn>
                Loading...
              </TableRowColumn>
            </TableRow> :
            pulls.data.map((pull, index) => {
              return <PullRow pull={pull} key={index} index={index}/>;
            })}
        </TableBody>
      </Table>
    );
  }
}

const mapStateToProps = (state) => {
  return {
    isFetching: state.reducer.isFetching,
    pulls: state.reducer.pulls,
    repoId: state.router.repoId
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    initialize: () => { dispatch(initialize()) }
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(PullsTable);
